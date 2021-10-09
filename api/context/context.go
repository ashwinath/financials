package context

import (
	"fmt"

	"github.com/ashwinath/financials/api/config"
	mediator "github.com/ashwinath/financials/api/mediators"
	"github.com/ashwinath/financials/api/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Context contains all the dependencies as part of DI
type Context struct {
	// Database
	DB *gorm.DB

	// Services
	UserService             *service.UserService
	SessionService          *service.SessionService
	TradeTransactionService *service.TradeService
	SymbolService           *service.SymbolService
	AlphaVantageService     *service.AlphaVantageService
	ExchangeRateService     *service.ExchangeRateService

	// Mediators
	LoginMediator *mediator.LoginMediator
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

	// Services
	context.SessionService = service.NewSessionService(db)
	context.UserService = service.NewUserService(db)
	context.TradeTransactionService = service.NewTradeService(db, c.Database.BatchInsertSize)
	context.SymbolService = service.NewSymbolService(db)
	context.AlphaVantageService = service.NewAlphaVantageService(c.AlphaVantage.APIKey)
	context.ExchangeRateService = service.NewExchangeRateService(db, c.Database.BatchInsertSize)

	// Mediators
	context.LoginMediator = mediator.NewLoginMediator(
		context.UserService,
		context.SessionService,
	)
	context.TradeMediator = mediator.NewTradeMediator(
		context.TradeTransactionService,
		context.SymbolService,
		context.AlphaVantageService,
		context.ExchangeRateService,
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
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
