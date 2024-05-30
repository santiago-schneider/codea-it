package models

type LtpDto struct {
	Ltp []LtpElement `json:"ltp"`
}

type LtpElement struct {
	Amount string `json:"amount"`
	Pair   string `json:"pair"`
}
