package mediator

import (
	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
)

// TradeMediator handles everything regarding trades
type TradeMediator struct {
	tradeTransactionService *service.TradeService
}

// NewTradeMediator creates a new NewTradeMediator
func NewTradeMediator(
	tradeTransactionService *service.TradeService,
) *TradeMediator {
	return &TradeMediator{
		tradeTransactionService: tradeTransactionService,
	}
}

// CreateTransactionInBulk creates multiple trade transactions at once
func (m *TradeMediator) CreateTransactionInBulk(
	session *models.Session,
	transactions []models.Trade,
) error {
	for _, tx := range transactions {
		tx.UserID = session.UserID
	}

	return m.tradeTransactionService.BulkAdd(transactions)
}
