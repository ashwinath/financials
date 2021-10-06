package mediator

import (
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
	tradeService *service.TradeService
}

// NewTradeMediator creates a new NewTradeMediator
func NewTradeMediator(
	tradeService *service.TradeService,
) *TradeMediator {
	return &TradeMediator{
		tradeService: tradeService,
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

	return m.tradeService.BulkAdd(transactions)
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
