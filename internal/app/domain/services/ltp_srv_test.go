package services

import (
	"errors"
	"testing"
	"time"

	"codea-it/internal/app/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExternalAPI es un mock de la interfaz ExternalAPI.
type MockExternalAPI struct {
	mock.Mock
}

func (m *MockExternalAPI) FetchData(currencyPair string) (models.TickerResponse, error) {
	args := m.Called(currencyPair)
	return args.Get(0).(models.TickerResponse), args.Error(1)
}

// MockRepository es un mock de la interfaz Repository.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetRecord(currencyPair string) (models.Record, error) {
	args := m.Called(currencyPair)
	return args.Get(0).(models.Record), args.Error(1)
}

func (m *MockRepository) SaveRecord(currencyPair string, expirationDate time.Time, status string) error {
	args := m.Called(currencyPair, expirationDate, status)
	return args.Error(0)
}

func (m *MockRepository) UpdateRecord(currencyPair string, expirationDate time.Time, status string) error {
	args := m.Called(currencyPair, expirationDate, status)
	return args.Error(0)
}

func (m *MockRepository) SavePairs(currencyPair string, pair string, amount string) error {
	args := m.Called(currencyPair, pair, amount)
	return args.Error(0)
}

func (m *MockRepository) UpdatePairs(currencyPair string, pair string, amount string) error {
	args := m.Called(currencyPair, pair, amount)
	return args.Error(0)
}

func (m *MockRepository) GetPairs(currencyPair string) ([]models.Pair, error) {
	args := m.Called(currencyPair)
	return args.Get(0).([]models.Pair), args.Error(1)
}

func TestExtractTickerData(t *testing.T) {

	t.Run("should return error if repository GetRecord fails", func(t *testing.T) {
		mockAPI := new(MockExternalAPI)
		mockRepo := new(MockRepository)
		service := NewLtpService(mockAPI, mockRepo)

		mockRepo.On("GetRecord", currencyPair).Return(models.Record{}, errors.New("db error"))

		_, err := service.ExtractTickerData()
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when data is not expired but data fetch previously fails", func(t *testing.T) {
		mockAPI := new(MockExternalAPI)
		mockRepo := new(MockRepository)
		service := NewLtpService(mockAPI, mockRepo)

		expiredRecord := models.Record{
			ExpirationDate: time.Now().Add(1 * time.Minute),
			Status:         failedStatus,
		}

		mockRepo.On("GetRecord", currencyPair).Return(expiredRecord, nil)

		_, err := service.ExtractTickerData()
		assert.Error(t, err)
		assert.Equal(t, "failed to get data", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if data is expired and fetchAndUpdateLtp fails", func(t *testing.T) {
		mockAPI := new(MockExternalAPI)
		mockRepo := new(MockRepository)
		service := NewLtpService(mockAPI, mockRepo)

		expiredRecord := models.Record{
			ExpirationDate: time.Now().Add(-1 * time.Minute),
			Status:         successStatus,
		}
		mockRepo.On("GetRecord", currencyPair).Return(expiredRecord, nil)
		mockAPI.On("FetchData", currencyPair).Return(models.TickerResponse{}, errors.New("api error"))
		mockRepo.On("UpdateRecord", currencyPair, mock.Anything, failedStatus).Return(nil)

		_, err := service.ExtractTickerData()
		assert.Error(t, err)
		assert.Equal(t, "failed to fetch data from external API", err.Error())
		mockRepo.AssertExpectations(t)
		mockAPI.AssertExpectations(t)
	})

	t.Run("should return data from repository if not expired", func(t *testing.T) {
		mockAPI := new(MockExternalAPI)
		mockRepo := new(MockRepository)
		service := NewLtpService(mockAPI, mockRepo)

		validRecord := models.Record{
			ExpirationDate: time.Now().Add(1 * time.Minute),
			Status:         successStatus,
		}
		mockRepo.On("GetRecord", currencyPair).Return(validRecord, nil)
		mockRepo.On("GetPairs", currencyPair).Return([]models.Pair{
			{Pair: "BTC/USD", Ltp: "30000"},
			{Pair: "BTC/CHF", Ltp: "28000"},
			{Pair: "BTC/EUR", Ltp: "27000"},
		}, nil)

		result, err := service.ExtractTickerData()
		assert.NoError(t, err)
		assert.Equal(t, 3, len(result.Ltp))
		mockRepo.AssertExpectations(t)
	})
}
