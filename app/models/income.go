package models

import "time"

// Income is the Income received
type Income struct {
	Model
	TransactionDate time.Time `json:"transaction_date"`
	Type            string    `json:"type"`
	// This number can be both negative or positive
	Amount float64 `json:"amount"`
}
