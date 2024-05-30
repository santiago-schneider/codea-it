package ports

import "codea-it/internal/app/domain/models"

type ExternalAPI interface {
	FetchData(param string) (models.TickerResponse, error)
}
