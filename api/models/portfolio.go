package models

import "time"

// Portfolio is the rate for a particular date
type Portfolio struct {
	Model
	UserID        string    `json:"-"`
	TradeDate     time.Time `json:"trade_date"`
	Symbol        string    `json:"symbol"`
	Quantity      float64   `json:"quantity"`
	Principal     float64   `json:"principal"`
	NAV           float64   `json:"nav" gorm:"column:nav"`
	SimpleReturns float64   `json:"simple_returns"`
}
