package mediator

import (
	"log"
	"strings"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
)

var (
	orderByDatePurchased = "date_purchased"
	orderDirection       = "desc"
)

// TradeMediator handles everything regarding trades
type TradeMediator struct {
	tradeService  *service.TradeService
	symbolService *service.SymbolService
}

// NewTradeMediator creates a new NewTradeMediator
func NewTradeMediator(
	tradeService *service.TradeService,
	symbolService *service.SymbolService,
) *TradeMediator {
	return &TradeMediator{
		tradeService:  tradeService,
		symbolService: symbolService,
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
