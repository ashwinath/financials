package models

import "time"

// Asset is the Asset or liability, negative value for liabilities
type Asset struct {
	Model
	TransactionDate time.Time `json:"transaction_date"`
	Type            string    `json:"type"`
	// This number can be both negative or positive
	Amount float64 `json:"amount"`
}
