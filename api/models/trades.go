package models

import "time"

// Trade is one transaction in the stock exchange.
type Trade struct {
	Model
	UserID        string    `json:"-"`
	DatePurchased time.Time `json:"date_purchased" validate:"required"`
	Symbol        string    `json:"symbol" validate:"required"`
	PriceEach     float64   `json:"price_each" validate:"required"`
	Quantity      float64   `json:"quantity" validate:"required"`
}
