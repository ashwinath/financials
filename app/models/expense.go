package models

import "time"

// Expense is the expense made
type Expense struct {
	Model
	TransactionDate time.Time `json:"transaction_date"`
	Type            string    `json:"type"`
	// This number can be both negative or positive
	// Postitive denotes amount spent
	// Negative denotes amount reimbursed
	Amount float64 `json:"amount"`
}
