package outbound

import (
	"codea-it/internal/app/domain/models"
	"codea-it/internal/app/ports"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type KrakenAPIClient struct {
	baseURL string
}

func NewKrakenAPIClient() ports.ExternalAPI {
	return &KrakenAPIClient{
		baseURL: "https://api.kraken.com",
	}
}

func (c *KrakenAPIClient) FetchData(param string) (models.TickerResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/0/public/Ticker?pair=%s", c.baseURL, param))
	if err != nil {
		return models.TickerResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body", err)
		}
	}(resp.Body)

	var result models.TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.TickerResponse{}, err
	}
	return result, nil
}
