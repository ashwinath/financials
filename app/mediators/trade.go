package mediator

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"github.com/ashwinath/financials/api/utils"
)

var (
	orderByDatePurchased = "date_purchased"
	orderDirection       = "desc"
	symbolQueryPageSize  = 100
	tradeQueryPageSize   = 1000
	tradeOrderBy         = "date_purchased"
	tradeOrder           = "asc"
)

const (
	baseCurrency  = "SGD"
	secondsInADay = 60 * 60 * 24
)

// TradeMediator handles everything regarding trades
type TradeMediator struct {
	tradeService        *service.TradeService
	symbolService       *service.SymbolService
	alphaVantageService *service.AlphaVantageService
	exchangeRateService *service.ExchangeRateService
	stockService        *service.StockService
	portfolioService    *service.PortfolioService
	csvPath             string
}

// NewTradeMediator creates a new NewTradeMediator
func NewTradeMediator(
	tradeService *service.TradeService,
	symbolService *service.SymbolService,
	alphaVantageService *service.AlphaVantageService,
	exchangeRateService *service.ExchangeRateService,
	stockService *service.StockService,
	portfolioService *service.PortfolioService,
	csvPath string,
) *TradeMediator {
	return &TradeMediator{
		tradeService:        tradeService,
		symbolService:       symbolService,
		alphaVantageService: alphaVantageService,
		exchangeRateService: exchangeRateService,
		stockService:        stockService,
		portfolioService:    portfolioService,
		csvPath:             csvPath,
	}
}

func (m *TradeMediator) insertTradesWithCSV() error {
	records, err := utils.ReadCSV(m.csvPath)
	if err != nil {
		return err
	}

	symbolSet := make(map[string]struct{})
	headers := records[0]
	var trades []*models.Trade
	for recordNum := 1; recordNum < len(records); recordNum++ {
		trade := &models.Trade{}
		for i, value := range records[recordNum] {
			switch headers[i] {
			case "date_purchased":
				layout := "2006-01-02T15:04:05.000Z"
				str := fmt.Sprintf("%sT08:00:00.000Z", value)
				t, err := time.Parse(layout, str)
				if err != nil {
					return err
				}
				trade.DatePurchased = t
			case "symbol":
				trade.Symbol = value
				symbolSet[value] = struct{}{}
			case "trade_type":
				trade.TradeType = value
			case "price_each":
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					trade.PriceEach = v
				} else {
					return err
				}
			case "quantity":
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					trade.Quantity = v
				} else {
					return err
				}
			}
		}
		trades = append(trades, trade)
	}

	// Create symbol set
	m.createSymbolIfNotExists(symbolSet)

	// Truncate the table first since we don't want to deal with duplicate trades
	if err := m.tradeService.TruncateTable(); err != nil {
		return err
	}

	return m.tradeService.BulkAdd(trades)
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

			if err != nil {
				return err
			}

			// Due to time zone differences, we just add one extra day, assume stocks are us time based
			t = t.Add(secondsInADay * time.Second)

			if maxDate == nil || maxDate.Before(t) {
				maxDate = &t
			}

			// TODO: Handle split coefficient.
			price, err := strconv.ParseFloat(value.Close, 64)
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

// Data might be dirty so we take the previous entry so that it's most accurate
func (m *TradeMediator) getCurrency(symbol string, date time.Time) *models.ExchangeRate {
	currencyDate := date
	er, err := m.exchangeRateService.Find(symbol, currencyDate)
	for err != nil {
		// Keep finding the last valid entry
		currencyDate = currencyDate.Add(-1 * secondsInADay * time.Second)
		er, err = m.exchangeRateService.Find(symbol, currencyDate)
	}
	return er
}

// Data might be dirty so we take the previous entry so that it's most accurate
func (m *TradeMediator) getStock(symbol string, date time.Time) *models.Stock {
	currencyDate := date
	stonk, err := m.stockService.Find(symbol, currencyDate)
	for err != nil {
		// Keep finding the last valid entry
		currencyDate = currencyDate.Add(-1 * secondsInADay * time.Second)
		stonk, err = m.stockService.Find(symbol, currencyDate)
	}
	return stonk
}

func (m *TradeMediator) calculatePortfolio() error {
	stockSymbolsSet := make(map[string]struct{})
	var allUserTrades []models.Trade

	// iterate through all user trades
	page := 0
	totalPages := 1 // just to satisfy initial condition
	for page < totalPages {
		paginatedTrades, err := m.tradeService.List(service.TradeListOptions{
			PaginationOptions: service.PaginationOptions{
				PageSize: &tradeQueryPageSize,
				OrderBy:  &tradeOrderBy,
				Order:    &tradeOrder,
			},
		})
		if err != nil {
			return err
		}

		trades := paginatedTrades.Results.([]models.Trade)
		for _, trade := range trades {
			stockSymbolsSet[trade.Symbol] = struct{}{}
		}
		allUserTrades = append(allUserTrades, trades...)

		totalPages = paginatedTrades.Paging.Pages
		page++
	}

	// Get base currencies required
	currencySymbolsSet := make(map[string]struct{})
	stockCurrencyMap := make(map[string]string)
	for stockSymbol := range stockSymbolsSet {
		symbol, err := m.symbolService.Find(stockSymbol)
		if err != nil {
			return err
		}

		currencySymbolsSet[symbol.BaseCurrency] = struct{}{}
		stockCurrencyMap[stockSymbol] = symbol.BaseCurrency
	}

	// Process portfolios partially for days with active trading
	allPortfoliosMap := make(map[string][]models.Portfolio)
	lastPortfolioMap := make(map[string]models.Portfolio)
	for _, trade := range allUserTrades {
		exchangeRate := m.getCurrency(stockCurrencyMap[trade.Symbol], trade.DatePurchased)

		tradeMultiplier := 1.0
		if trade.TradeType == "sell" {
			tradeMultiplier = -1.0
		}
		var portfolio models.Portfolio
		if lastPortfolio, ok := lastPortfolioMap[trade.Symbol]; !ok {
			// Base case
			principal := trade.PriceEach * trade.Quantity * exchangeRate.Price
			portfolio = models.Portfolio{
				TradeDate: trade.DatePurchased,
				Symbol:    trade.Symbol,
				Principal: principal,
				Quantity:  trade.Quantity,
			}
		} else {
			// any other increasing trade.
			principal := lastPortfolio.Principal + trade.PriceEach*trade.Quantity*exchangeRate.Price*tradeMultiplier
			portfolio = models.Portfolio{
				TradeDate: trade.DatePurchased,
				Symbol:    trade.Symbol,
				Principal: principal,
				Quantity:  lastPortfolio.Quantity + (trade.Quantity * tradeMultiplier),
			}
		}

		lastPortfolioMap[trade.Symbol] = portfolio
		allPortfoliosMap[trade.Symbol] = append(allPortfoliosMap[trade.Symbol], portfolio)
	}

	var allPortfolios []models.Portfolio
	// Another pass through to fill gaps so that the data is continuous daily
	loc, _ := time.LoadLocation("Asia/Singapore")
	for symbol, partialPortfolios := range allPortfoliosMap {
		currentDate := partialPortfolios[0].TradeDate
		portfolioWithTradesMap := make(map[time.Time]models.Portfolio)
		for _, portfolio := range partialPortfolios {
			// There might be multiple trades in a single day for each symbol, we need to combine them
			if existingPortfolio, ok := portfolioWithTradesMap[portfolio.TradeDate]; ok {
				portfolio.Quantity = existingPortfolio.Quantity
				portfolio.Principal = existingPortfolio.Principal
			}
			portfolioWithTradesMap[portfolio.TradeDate] = portfolio
		}

		now := time.Now()
		tomorrow := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			16, 0, 0, 0, loc,
		)
		var previousPortfolio models.Portfolio
		for currentDate.Before(tomorrow) {
			er := m.getCurrency(stockCurrencyMap[symbol], currentDate)

			exchangeRate := er.Price
			newPortfolio := previousPortfolio
			tradeDate := currentDate
			newPortfolio.TradeDate = tradeDate
			newPortfolio.Symbol = symbol

			if tradePortfolio, ok := portfolioWithTradesMap[currentDate]; ok {
				// new trades here, we need to merge old value with new value
				newPortfolio.Quantity = tradePortfolio.Quantity
				newPortfolio.Principal = tradePortfolio.Principal
			}

			// update every nav and simple returns
			stock := m.getStock(symbol, currentDate)
			newPortfolio.NAV = newPortfolio.Quantity * stock.Price * exchangeRate
			newPortfolio.SimpleReturns = (newPortfolio.NAV - newPortfolio.Principal) / newPortfolio.Principal
			allPortfolios = append(allPortfolios, newPortfolio)

			currentDate = currentDate.Add(secondsInADay * time.Second)
			previousPortfolio = newPortfolio
		}
	}

	err := m.portfolioService.BulkAdd(allPortfolios)
	if err != nil {
		return err
	}

	return nil
}

// ProcessTrades processes the trades.
func (m *TradeMediator) ProcessTrades() error {
	start := time.Now()
	log.Printf("Running one round of process trades")

	// Get trades from csv
	err := m.insertTradesWithCSV()
	if err != nil {
		return fmt.Errorf("error parsing csv: %s", err)
	}

	// Gets currencies involved and synchronises the tables
	m.syncSymbolTable()

	// Gets the exchange rates for all currencies to SGD
	err = m.processCurrency()
	if err != nil {
		return fmt.Errorf("error downloading currency information: %s", err)
	}

	// Query stock rates and put into table
	err = m.processStocks()
	if err != nil {
		return fmt.Errorf("error downloading stocks information: %s", err)
	}

	// Update portfolio
	err = m.calculatePortfolio()
	if err != nil {
		return fmt.Errorf("error calculating portfolio information: %s", err)
	}

	log.Printf("Finished one round of process trades, time taken: %s", time.Since(start))
	return nil
}
