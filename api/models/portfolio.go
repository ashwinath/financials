package models

import "time"

// Portfolio is the rate for a particular date
type Portfolio struct {
	Model
	UserID        string
	TradeDate     time.Time
	Symbol        string
	Quantity      float64
	Principal     float64
	NAV           float64 `gorm:"column:nav"`
	SimpleReturns float64
}
