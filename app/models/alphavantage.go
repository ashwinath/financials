package models

// AlphaVantageBestMatches is the result when you query the following
// https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=<symbol>&apikey=<apikey>
type AlphaVantageBestMatches struct {
	BestMatches []AlphaVantageSymbolSearchResult `json:"bestMatches"`
}

// AlphaVantageSymbolSearchResult is the search result item in the list
type AlphaVantageSymbolSearchResult struct {
	Symbol   string `json:"2. symbol"`
	Currency string `json:"8. currency"`
}

// AlphaVantageCurrencyResult contains a map of daily forex values
type AlphaVantageCurrencyResult struct {
	Results map[string]AlphaVantageCurrencyDailyResult `json:"Time Series FX (Daily)"`
}

// AlphaVantageCurrencyDailyResult contains the single value of a daily currency result
type AlphaVantageCurrencyDailyResult struct {
	Close string `json:"4. close"`
}

// AlphaVantageStockResult contains a map of daily stock values
type AlphaVantageStockResult struct {
	Results map[string]AlphaVantageStockDailyResult `json:"Time Series (Daily)"`
}

// AlphaVantageStockDailyResult contains the single value of a daily stock result
type AlphaVantageStockDailyResult struct {
	// Alphavantage just made this a premium feature so we have to fix splits next time manually
	AdjustedClose string `json:"5. adjusted close"`
	// Using close after adjusted became a premium feature
	Close string `json:"4. close"`
}
