package mediator

import (
	"testing"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionInBulk(t *testing.T) {
	db, err := service.CreateTestDB()
	assert.Nil(t, err)

	svc := service.NewTradeService(db, 20)

	session := &models.Session{
		UserID: "hello",
	}

	m := NewTradeMediator(svc)

	trades := models.CreateTradeTransactions(5)
	err = m.CreateTransactionInBulk(session, trades)
	for _, trade := range trades {
		defer svc.Delete(&trade)
	}

	for _, trade := range trades {
		tr, err := svc.Find(trade.ID)
		assert.Nil(t, err)
		assert.NotNil(t, tr)

		assert.NotEqual(t, "", tr.UserID)
	}
}
