package mediator

import (
	"testing"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTransactionInBulk(t *testing.T) {
	service.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
		tradeSvc := service.NewTradeService(db, 20)
		symbolSvc := service.NewSymbolService(db)
		m := TradeMediator{
			tradeService:  tradeSvc,
			symbolService: symbolSvc,
		}

		session := &models.Session{
			UserID: "hello",
		}
		trades := models.CreateTradeTransactions(5)
		err := m.CreateTransactionInBulk(session, trades)
		assert.Nil(t, err)

		for _, trade := range trades {
			tr, err := tradeSvc.Find(trade.ID)
			assert.Nil(t, err)
			assert.NotNil(t, tr)

			assert.NotEqual(t, "", tr.UserID)
		}
	})
}
