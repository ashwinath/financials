package models

import "time"

// ExchangeRate is the rate for a particular date
type ExchangeRate struct {
	Model
	TradeDate *time.Time
	Symbol    string
	Price     float64
}
