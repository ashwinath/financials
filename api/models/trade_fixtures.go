package models

import "time"

// CreateTradeTransactions is a helper function to create dummy fixtures
// For Test purposes only
func CreateTradeTransactions(count int) []*Trade {
	trades := []*Trade{}
	for i := 0; i < count; i++ {
		trades = append(trades, &Trade{
			DatePurchased: time.Now().Add(time.Minute * time.Duration(count)),
			Symbol:        "VWRA.LON",
			PriceEach:     100.25,
			Quantity:      100,
		})
	}
	return trades
}
