package service

import (
	"fmt"

	"github.com/ashwinath/financials/api/models"
)

const (
	searchURLFormat = "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=%s&apikey=%s"
	fxURLFormat     = "https://www.alphavantage.co/query?function=FX_DAILY&from_symbol=%s&to_symbol=SGD&apikey=%s&outputsize=%s"
)

// AlphaVantageService is an external service that queries stock info
type AlphaVantageService struct {
	apiKey string
}

// NewAlphaVantageService a new SessionService
func NewAlphaVantageService(apiKey string) *AlphaVantageService {
	return &AlphaVantageService{
		apiKey: apiKey,
	}
}

// GetStockInfo gets the stock information of a stock.
func (s *AlphaVantageService) GetStockInfo(symbol string) (*models.AlphaVantageSymbolSearchResult, error) {
	var result models.AlphaVantageBestMatches
	err := query(
		fmt.Sprintf(searchURLFormat, symbol, s.apiKey),
		&result,
	)
	if err != nil {
		return nil, err
	}

	if len(result.BestMatches) == 0 {
		return nil, fmt.Errorf("no such symbol")
	}

	return &result.BestMatches[0], nil
}

// GetCurrencyHistory gets the currency history
func (s *AlphaVantageService) GetCurrencyHistory(symbol string, isCompact bool) (*models.AlphaVantageCurrencyResult, error) {
	var result models.AlphaVantageCurrencyResult
	outputSize := "full"
	if isCompact {
		outputSize = "compact"
	}
	err := query(
		fmt.Sprintf(fxURLFormat, symbol, s.apiKey, outputSize),
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
