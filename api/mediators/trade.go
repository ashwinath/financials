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
	portfolioCalculationInterval time.Duration
	tradeService                 *service.TradeService
	symbolService                *service.SymbolService
	alphaVantageService          *service.AlphaVantageService
	exchangeRateService          *service.ExchangeRateService
	stockService                 *service.StockService
	userService                  *service.UserService
	portfolioService             *service.PortfolioService
}

// NewTradeMediator creates a new NewTradeMediator
func NewTradeMediator(
	portfolioCalculationInterval time.Duration,
	tradeService *service.TradeService,
	symbolService *service.SymbolService,
	alphaVantageService *service.AlphaVantageService,
	exchangeRateService *service.ExchangeRateService,
	stockService *service.StockService,
	userService *service.UserService,
	portfolioService *service.PortfolioService,
) *TradeMediator {
	return &TradeMediator{
		portfolioCalculationInterval: portfolioCalculationInterval,
		tradeService:                 tradeService,
		symbolService:                symbolService,
		alphaVantageService:          alphaVantageService,
		exchangeRateService:          exchangeRateService,
		stockService:                 stockService,
		userService:                  userService,
		portfolioService:             portfolioService,
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

			if err != nil {
				return err
			}

			// Due to time zone differences, we just add one extra day, assume stocks are us time based
			t = t.Add(secondsInADay * time.Second)

			if maxDate == nil || maxDate.Before(t) {
				maxDate = &t
			}

			// TODO: Handle split coefficient.
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

func (m *TradeMediator) calculatePortfolio() error {
	users, err := m.userService.ListUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		stockSymbolsSet := make(map[string]struct{})
		var allUserTrades []models.Trade

		// iterate through all user trades
		page := 0
		totalPages := 1 // just to satisfy initial condition
		for page < totalPages {
			paginatedTrades, err := m.tradeService.List(service.TradeListOptions{
				UserID: &user.ID,
				PaginationOptions: service.PaginationOptions{
					PageSize: &tradeQueryPageSize,
					OrderBy:  &tradeOrderBy,
					Order:    &tradeOrder,
				},
			})
			if err != nil {
				return nil
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
				return nil
			}

			currencySymbolsSet[symbol.BaseCurrency] = struct{}{}
			stockCurrencyMap[stockSymbol] = symbol.BaseCurrency
		}

		// Process portfolios partially for days with active trading
		allPortfoliosMap := make(map[string][]models.Portfolio)
		lastPortfolioMap := make(map[string]models.Portfolio)
		for _, trade := range allUserTrades {
			exchangeRate := m.getCurrency(stockCurrencyMap[trade.Symbol], trade.DatePurchased)

			var portfolio models.Portfolio
			if lastPortfolio, ok := lastPortfolioMap[trade.Symbol]; !ok {
				// Base case
				principal := trade.PriceEach * trade.Quantity * exchangeRate.Price
				portfolio = models.Portfolio{
					UserID:    trade.UserID,
					TradeDate: trade.DatePurchased,
					Symbol:    trade.Symbol,
					Principal: principal,
					Quantity:  trade.Quantity,
				}
			} else {
				// any other increasing trade.
				tradeMultiplier := 1.0
				if trade.TradeType == "sell" {
					tradeMultiplier = -1.0
				}
				principal := lastPortfolio.Principal + trade.PriceEach*trade.Quantity*exchangeRate.Price*tradeMultiplier
				portfolio = models.Portfolio{
					UserID:    trade.UserID,
					TradeDate: trade.DatePurchased,
					Symbol:    trade.Symbol,
					Principal: principal,
					Quantity:  lastPortfolio.Quantity + trade.Quantity*tradeMultiplier,
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
					portfolio.Quantity += existingPortfolio.Quantity
					portfolio.Principal += existingPortfolio.Principal
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
				newPortfolio.UserID = user.ID

				if tradePortfolio, ok := portfolioWithTradesMap[currentDate]; ok {
					// new trades here, we need to merge old value with new value
					newPortfolio.Quantity = tradePortfolio.Quantity
					newPortfolio.Principal = tradePortfolio.Principal
				}

				if stock, err := m.stockService.Find(symbol, currentDate); err == nil {
					// Trading day, update NAV and simple returns
					newPortfolio.NAV = newPortfolio.Quantity * stock.Price * exchangeRate
					newPortfolio.SimpleReturns = (newPortfolio.NAV - newPortfolio.Principal) / newPortfolio.Principal
				}
				allPortfolios = append(allPortfolios, newPortfolio)

				currentDate = currentDate.Add(secondsInADay * time.Second)
				previousPortfolio = newPortfolio
			}
		}
		err = m.portfolioService.BulkAdd(allPortfolios)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProcessTrades processes all the trades for all users.
// Since I'm the only user it's not going to be optimised to calculate parallely
func (m *TradeMediator) ProcessTrades() {
	for {
		start := time.Now()
		log.Printf("Running one round of process trades")
		// Gets currencies involved and synchronises the tables
		m.syncSymbolTable()

		// Gets the exchange rates for all currencies to SGD
		err := m.processCurrency()
		if err != nil {
			log.Printf("error downloading currency information: %s", err)
			goto endloop
		}

		// Query stock rates and put into table
		err = m.processStocks()
		if err != nil {
			log.Printf("error downloading stocks information: %s", err)
			goto endloop
		}

		// Update portfolio for each user id
		err = m.calculatePortfolio()
		if err != nil {
			log.Printf("error calculating portfolio information: %s", err)
			goto endloop
		}
		log.Printf("Finished one round of process trades, time taken: %s", time.Since(start))

	endloop:
		log.Printf("ProcessTrades Sleeping for: %s", m.portfolioCalculationInterval)
		time.Sleep(m.portfolioCalculationInterval)
	}
}
