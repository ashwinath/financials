package models

import "time"

// Symbol contains the stock/exchange rate symbol
type Symbol struct {
	Model
	SymbolType        SymbolType
	Symbol            string
	BaseCurrency      string
	LastProcessedDate *time.Time
}

// SymbolType is a type of symbol
type SymbolType string

const (
	// SymbolStock is a stock type
	SymbolStock = SymbolType("stock")
	// SymbolCurrency is a currency type
	SymbolCurrency = SymbolType("currency")
)
