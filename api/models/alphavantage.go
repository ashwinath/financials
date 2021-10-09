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
