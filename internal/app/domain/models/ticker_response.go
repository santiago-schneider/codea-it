package models

type TickerResponse struct {
	Result Result        `json:"result"`
	Error  []interface{} `json:"error"`
}

type Result struct {
	Xbtchf   Data `json:"XBTCHF"`
	Xxbtzusd Data `json:"XXBTZUSD"`
	Xxbtzeur Data `json:"XXBTZEUR"`
}

type Data struct {
	Price   []string `json:"p"`
	Ask     []string `json:"a"`
	Bid     []string `json:"b"`
	Closed  []string `json:"c"`
	Trades  []int64  `json:"t"`
	Volume  []string `json:"v"`
	High    []string `json:"h"`
	Low     []string `json:"l"`
	Opening string   `json:"o"`
}
