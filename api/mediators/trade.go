package mediator

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
)

var (
	orderByDatePurchased = "date_purchased"
	orderDirection       = "desc"
	symbolQueryPageSize  = 100
)

const (
	baseCurrency = "SGD"
)

// TradeMediator handles everything regarding trades
type TradeMediator struct {
	tradeService        *service.TradeService
	symbolService       *service.SymbolService
	alphaVantageService *service.AlphaVantageService
	exchangeRateService *service.ExchangeRateService
	stockService        *service.StockService
}

// NewTradeMediator creates a new NewTradeMediator
func NewTradeMediator(
	tradeService *service.TradeService,
	symbolService *service.SymbolService,
	alphaVantageService *service.AlphaVantageService,
	exchangeRateService *service.ExchangeRateService,
	stockService *service.StockService,
) *TradeMediator {
	return &TradeMediator{
		tradeService:        tradeService,
		symbolService:       symbolService,
		alphaVantageService: alphaVantageService,
		exchangeRateService: exchangeRateService,
		stockService:        stockService,
	}
}

// CreateTransactionInBulk creates multiple trade transactions at once
func (m *TradeMediator) CreateTransactionInBulk(
	session *models.Session,
	transactions []*models.Trade,
) error {
	for _, tx := range transactions {
		tx.UserID = session.UserID
		tx.Symbol = strings.ToUpper(tx.Symbol)
	}

	symbolSet := make(map[string]struct{})
	for _, tx := range transactions {
		symbolSet[tx.Symbol] = struct{}{}
	}

	go m.createSymbolIfNotExists(symbolSet)

	return m.tradeService.BulkAdd(transactions)
}

func (m *TradeMediator) createSymbolIfNotExists(symbolSet map[string]struct{}) {
	// If it's slow then we can optimise this later
	for key := range symbolSet {
		_, err := m.symbolService.Find(key)
		if err == nil {
			// It exists, skip
			continue
		}

		// Did not find
		err = m.symbolService.Save(&models.Symbol{
			SymbolType: models.SymbolStock,
			Symbol:     key,
		})

		if err != nil {
			log.Printf("Could not insert symbol: %s", err.Error())
		}
	}
}

// ListTrades lists all the trades
func (m *TradeMediator) ListTrades(
	session *models.Session,
	options service.TradeListOptions,
) (*service.PaginatedResults, error) {
	if options.OrderBy == nil {
		options.OrderBy = &orderByDatePurchased
	}

	if options.Order == nil {
		options.Order = &orderDirection
	}

	options.UserID = &session.UserID

	return m.tradeService.List(options)
}

func (m *TradeMediator) syncSymbolTable() {
	// Query symbol table
	symbolsPaginated, err := m.symbolService.List(service.SymbolListOptions{
		PaginationOptions: service.PaginationOptions{
			PageSize: &symbolQueryPageSize,
		},
		SymbolType: models.SymbolStock,
	})

	if err != nil {
		log.Printf("Error querying symbol table: %s", err.Error())
	}
	// We are going to ignore the paginated results for now,
	// I don't expect more than 100 queries here since it will break the api limit also
	// I also am ignoring bulk insert as I don't think I will hold a lot of different etfs/currencies
	symbols := symbolsPaginated.Results.([]models.Symbol)
	for _, symbol := range symbols {
		if symbol.BaseCurrency == "" {
			// Query alpha vantage and populate the base currency
			result, err := m.alphaVantageService.GetStockInfo(symbol.Symbol)
			if err != nil {
				log.Printf("Error searching alpha vantage: %s", err.Error())
				continue
			}

			symbol.BaseCurrency = result.Currency
			err = m.symbolService.Save(&symbol)
			if err != nil {
				log.Printf("Error saving symbol base currency: %s", err.Error())
				continue
			}

		}
		currencySymbol := fmt.Sprintf("%s", symbol.BaseCurrency)
		_, err = m.symbolService.Find(currencySymbol)
		if err == nil {
			// Symbol exists, don't have to add it in
			continue
		}

		// else insert
		err = m.symbolService.Save(&models.Symbol{
			SymbolType: models.SymbolCurrency,
			Symbol:     currencySymbol,
		})
		if err != nil {
			log.Printf("Something went wrong inserting into symbols: %s", err)
		}
	}
}

func (m *TradeMediator) processCurrency() error {
	currencyPaginated, err := m.symbolService.List(service.SymbolListOptions{
		PaginationOptions: service.PaginationOptions{
			PageSize: &symbolQueryPageSize,
		},
		SymbolType: models.SymbolCurrency,
	})
	if err != nil {
		log.Printf("Error querying symbol table: %s", err.Error())
	}

	symbols := currencyPaginated.Results.([]models.Symbol)
	for _, symbol := range symbols {
		isCompact := true
		if symbol.LastProcessedDate == nil {
			isCompact = false
		}
		result, err := m.alphaVantageService.GetCurrencyHistory(symbol.Symbol, isCompact)
		if err != nil {
			return err
		}

		var allExchangeRates []*models.ExchangeRate
		var maxDate *time.Time
		for key, value := range result.Results {
			t, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT08:00:00.000Z", key))
			if maxDate == nil || maxDate.Before(t) {
				maxDate = &t
			}

			if err != nil {
				return err
			}

			price, err := strconv.ParseFloat(value.Close, 64)
			if err != nil {
				return err
			}

			allExchangeRates = append(allExchangeRates, &models.ExchangeRate{
				TradeDate: &t,
				Symbol:    symbol.Symbol,
				Price:     price,
			})
		}

		err = m.exchangeRateService.BulkAdd(allExchangeRates)
		if err != nil {
			return err
		}

		symbol.LastProcessedDate = maxDate
		err = m.symbolService.Save(&symbol)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *TradeMediator) processStocks() error {
	currencyPaginated, err := m.symbolService.List(service.SymbolListOptions{
		PaginationOptions: service.PaginationOptions{
			PageSize: &symbolQueryPageSize,
		},
		SymbolType: models.SymbolStock,
	})
	if err != nil {
		log.Printf("Error querying symbol table: %s", err.Error())
	}

	symbols := currencyPaginated.Results.([]models.Symbol)
	for _, symbol := range symbols {
		isCompact := true
		if symbol.LastProcessedDate == nil {
			isCompact = false
		}
		result, err := m.alphaVantageService.GetStockHistory(symbol.Symbol, isCompact)
		if err != nil {
			// Some fix required on this as this might stall.
			return err
		}

		var allStocks []*models.Stock
		var maxDate *time.Time
		for key, value := range result.Results {
			t, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT08:00:00.000Z", key))
			if maxDate == nil || maxDate.Before(t) {
				maxDate = &t
			}

			if err != nil {
				return err
			}

			price, err := strconv.ParseFloat(value.AdjustedClose, 64)
			if err != nil {
				return err
			}

			allStocks = append(allStocks, &models.Stock{
				TradeDate: &t,
				Symbol:    symbol.Symbol,
				Price:     price,
			})
		}

		err = m.stockService.BulkAdd(allStocks)
		if err != nil {
			return err
		}

		symbol.LastProcessedDate = maxDate
		err = m.symbolService.Save(&symbol)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessTrades processes all the trades for all users.
// Since I'm the only user it's not going to be optimised to calculate parallely
func (m *TradeMediator) ProcessTrades() {
	// Gets currencies involved and synchronises the tables
	m.syncSymbolTable()

	// Gets the exchange rates for all currencies to SGD
	err := m.processCurrency()
	if err != nil {
		log.Printf("error downloading currency information: %s", err)
		return
	}

	// Query stock rates and put into table
	err = m.processStocks()
	if err != nil {
		log.Printf("error downloading stocks information: %s", err)
		return
	}
	// Update portfolio for each user id
}
