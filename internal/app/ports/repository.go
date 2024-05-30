package ports

import (
	"codea-it/internal/app/domain/models"
	"time"
)

type Repository interface {
	SaveRecord(string, time.Time, string) error
	GetRecord(string) (models.Record, error)
	UpdateRecord(string, time.Time, string) error
	SavePairs(string, string, string) error
	GetPairs(string) ([]models.Pair, error)
	UpdatePairs(string, string, string) error
}
