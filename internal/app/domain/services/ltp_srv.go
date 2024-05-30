package services

import (
	"codea-it/internal/app/domain/models"
	"codea-it/internal/app/ports"
	"errors"
	"time"
)

type LtpService interface {
	ExtractTickerData() (models.LtpDto, error)
}

type ltpService struct {
	externalAPI ports.ExternalAPI
	repository  ports.Repository
}

func NewLtpService(apiClient ports.ExternalAPI, repository ports.Repository) LtpService {
	return &ltpService{
		externalAPI: apiClient,
		repository:  repository,
	}
}

const currencyPair = "BTCUSD,BTCCHF,BTCEUR"
const failedStatus = "FAILED"
const successStatus = "SUCCESS"

func (s *ltpService) ExtractTickerData() (models.LtpDto, error) {
	dbData, err := s.repository.GetRecord(currencyPair)
	if err != nil {
		return models.LtpDto{}, err
	}

	if dbData.ExpirationDate.After(time.Now()) && dbData.Status == failedStatus {
		return models.LtpDto{}, errors.New("failed to get data")
	}

	needsUpdate := dbData.ExpirationDate.Before(time.Now())
	expireDate := time.Now().Add(time.Minute * 1)
	var ltp models.LtpDto

	if needsUpdate {
		ltp, err = s.fetchAndUpdateLtp(currencyPair, expireDate, dbData)
		if err != nil {
			return models.LtpDto{}, err
		}
	} else {
		ltp, err = s.getLtp(currencyPair)
		if err != nil {
			return models.LtpDto{}, err
		}
	}

	return ltp, nil
}

func (s *ltpService) fetchAndUpdateLtp(currencyPair string, expireDate time.Time, dbData models.Record) (models.LtpDto, error) {
	data, err := s.externalAPI.FetchData(currencyPair)
	if err != nil || len(data.Error) > 0 {
		s.handleFailedUpdate(currencyPair, expireDate, dbData)
		return models.LtpDto{}, errors.New("failed to fetch data from external API")
	}

	ltp := models.LtpDto{
		Ltp: []models.LtpElement{
			{Amount: data.Result.Xbtchf.Closed[0], Pair: "BTC/CHF"},
			{Amount: data.Result.Xxbtzeur.Closed[0], Pair: "BTC/EUR"},
			{Amount: data.Result.Xxbtzusd.Closed[0], Pair: "BTC/USD"},
		},
	}

	var emptyRecord models.Record
	if dbData == emptyRecord {
		err = s.repository.SaveRecord(currencyPair, expireDate, successStatus)
		if err != nil {
			return models.LtpDto{}, err
		}

		for _, element := range ltp.Ltp {
			err = s.repository.SavePairs(currencyPair, element.Pair, element.Amount)
			if err != nil {
				return models.LtpDto{}, err
			}
		}
		return ltp, nil
	}

	err = s.repository.UpdateRecord(currencyPair, expireDate, successStatus)
	if err != nil {
		return models.LtpDto{}, err
	}

	for _, element := range ltp.Ltp {
		err = s.repository.UpdatePairs(currencyPair, element.Pair, element.Amount)
		if err != nil {
			return models.LtpDto{}, err
		}
	}

	return ltp, nil
}

func (s *ltpService) handleFailedUpdate(currencyPair string, expireDate time.Time, dbData models.Record) {
	if dbData == (models.Record{}) {
		_ = s.repository.SaveRecord(currencyPair, expireDate, failedStatus)
	} else {
		_ = s.repository.UpdateRecord(currencyPair, expireDate, failedStatus)
	}
}

func (s *ltpService) getLtp(currencyPair string) (models.LtpDto, error) {
	pairs, err := s.repository.GetPairs(currencyPair)
	if err != nil {
		return models.LtpDto{}, err
	}

	var ltp models.LtpDto
	for _, pair := range pairs {
		ltp.Ltp = append(ltp.Ltp, models.LtpElement{
			Amount: pair.Ltp,
			Pair:   pair.Pair,
		})
	}

	return ltp, nil
}
