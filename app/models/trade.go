package models

import "time"

// Trade is one transaction in the stock exchange.
type Trade struct {
	Model
	DatePurchased time.Time `json:"date_purchased" validate:"required"`
	Symbol        string    `json:"symbol" validate:"required"`
	TradeType     string    `json:"trade_type" validate:"required,oneof=buy sell"`
	PriceEach     float64   `json:"price_each" validate:"required"`
	Quantity      float64   `json:"quantity" validate:"required"`
}
