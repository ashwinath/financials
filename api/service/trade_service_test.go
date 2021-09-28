package service

import (
	"testing"

	"github.com/ashwinath/financials/api/models"
	"github.com/stretchr/testify/assert"
)

func TestBulkAdd(t *testing.T) {
	db, err := CreateTestDB()
	assert.Nil(t, err)

	svc := NewTradeService(db, 10)

	trades := models.CreateTradeTransactions(5)
	err = svc.BulkAdd(trades)
	assert.Nil(t, err)
	for _, trade := range trades {
		defer svc.Delete(&trade)
	}

	for _, trade := range trades {
		tr, err := svc.Find(trade.ID)
		assert.Nil(t, err)
		assert.NotNil(t, tr)
	}
}
