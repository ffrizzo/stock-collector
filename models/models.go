package models

import "time"

//StockHistory represents the history of Stock
type Stock struct {
	ID      int64      `json:"id"`
	Sell    float64    `json:"sell"`
	Rate    float64    `json:"rate"`
	Buy     float64    `json:"buy"`
	Ticker  string     `json:"ticker"`
	Account string     `json:"account"`
	User    string     `json:"user"`
	Time    *time.Time `json:"time"`
}
