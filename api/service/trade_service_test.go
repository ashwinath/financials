package service

import (
	"testing"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBulkAdd(t *testing.T) {
	t.Run("success | test bulk add", func(t *testing.T) {
		WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
			svc := NewTradeService(db, 10)

			trades := models.CreateTradeTransactions(5)
			err := svc.BulkAdd(trades)
			assert.Nil(t, err)

			for _, trade := range trades {
				tr, err := svc.Find(trade.ID)
				assert.Nil(t, err)
				assert.NotNil(t, tr)
			}
		})
	})
}

func TestList(t *testing.T) {
	userID := "hello-world"
	tests := []struct {
		name                 string
		userID               string
		options              TradeListOptions
		expectedLength       int
		hasError             bool
		expectedTotalRecords int64
		expectedPageNumber   int
		expectedPageCount    int
	}{
		{
			name: "success | nominal",
			options: TradeListOptions{
				UserID: &userID,
				PaginationOptions: PaginationOptions{
					Page:     nil,
					PageSize: nil,
					OrderBy:  nil,
					Order:    nil,
				},
			},
			expectedLength:       10,
			hasError:             false,
			expectedTotalRecords: 10,
			expectedPageNumber:   1,
			expectedPageCount:    1,
		},
		{
			name: "success | nominal",
			options: TradeListOptions{
				UserID: &userID,
				PaginationOptions: PaginationOptions{
					Page:     nil,
					PageSize: utils.Intref(5),
					OrderBy:  nil,
					Order:    nil,
				},
			},
			expectedLength:       5,
			hasError:             false,
			expectedTotalRecords: 10,
			expectedPageNumber:   1,
			expectedPageCount:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				insertedLength := 10
				svc := NewTradeService(db, insertedLength)

				trades := models.CreateTradeTransactions(10)
				for _, trade := range trades {
					trade.UserID = userID
				}

				err := svc.BulkAdd(trades)
				assert.Nil(t, err)

				result, err := svc.List(tt.options)
				if tt.hasError {
					assert.NotNil(t, err)
					return
				}

				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Results, tt.expectedLength)
				assert.Equal(t, result.Paging.Total, tt.expectedTotalRecords)
				assert.Equal(t, result.Paging.Page, tt.expectedPageNumber)
				assert.Equal(t, result.Paging.Pages, tt.expectedPageCount)
			})
		})
	}
}
