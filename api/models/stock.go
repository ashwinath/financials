package models

import "time"

// Stock is the rate for a particular date
type Stock struct {
	Model
	TradeDate *time.Time
	Symbol    string
	Price     float64
}
