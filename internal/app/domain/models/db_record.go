package models

import "time"

type Record struct {
	PairsNames     string    `json:"pairs_names,omitempty"`
	ExpirationDate time.Time `json:"expiration_date,omitempty"`
	Status         string    `json:"status,omitempty"`
}

type Pair struct {
	PairsNames string `json:"pairs_names,omitempty"`
	Pair       string `json:"pair,omitempty"`
	Ltp        string `json:"ltp,omitempty"`
}
