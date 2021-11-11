package context

import (
	"fmt"

	"github.com/ashwinath/financials/api/config"
	mediator "github.com/ashwinath/financials/api/mediators"
	"github.com/ashwinath/financials/api/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Context contains all the dependencies as part of DI
type Context struct {
	// All the configurations
	Config *config.Config

	// Database
	DB *gorm.DB

	// Services
	TradeTransactionService *service.TradeService
	SymbolService           *service.SymbolService
	AlphaVantageService     *service.AlphaVantageService
	ExchangeRateService     *service.ExchangeRateService
	StockService            *service.StockService
	PortfolioService        *service.PortfolioService

	// Mediators
	TradeMediator *mediator.TradeMediator
}

// InitContext inits all dependencies required by API server
func InitContext(c *config.Config) (*Context, error) {
	context := Context{}
	db, err := initDB(c.Database)
	if err != nil {
		return nil, err
	}
	context.DB = db
	context.Config = c

	// Services
	context.TradeTransactionService = service.NewTradeService(db, c.Database.BatchInsertSize)
	context.SymbolService = service.NewSymbolService(db)
	context.AlphaVantageService = service.NewAlphaVantageService(c.AlphaVantageAPIKey)
	context.ExchangeRateService = service.NewExchangeRateService(db, c.Database.BatchInsertSize)
	context.StockService = service.NewStockService(db, c.Database.BatchInsertSize)
	context.PortfolioService = service.NewPortfolioService(db, c.Database.BatchInsertSize)

	// Mediators
	context.TradeMediator = mediator.NewTradeMediator(
		context.TradeTransactionService,
		context.SymbolService,
		context.AlphaVantageService,
		context.ExchangeRateService,
		context.StockService,
		context.PortfolioService,
		c.TradesCSVFile,
	)

	return &context, nil
}

func initDB(dbConfig config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
		dbConfig.TimeZone,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
