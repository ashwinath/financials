package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
)

// TradeService is the interface to the database for the sessions tables
type TradeService struct {
	db              *gorm.DB
	batchInsertSize int
}

// NewTradeService creates a new UserService
func NewTradeService(db *gorm.DB, batchInsertSize int) *TradeService {
	return &TradeService{
		db:              db,
		batchInsertSize: batchInsertSize,
	}
}

// BulkAdd adds multiple transactions at once
func (s *TradeService) BulkAdd(transactions []*models.Trade) error {
	return s.db.CreateInBatches(transactions, s.batchInsertSize).Error
}

// TradeListOptions lists the trades from a user.
type TradeListOptions struct {
	PaginationOptions
}

// List lists all trades
func (s *TradeService) List(options TradeListOptions) (*PaginatedResults, error) {
	done := make(chan struct{}, 1)

	var count int64
	go func() {
		s.db.Model(&models.Trade{}).Count(&count)
		done <- struct{}{}
	}()

	var results []models.Trade
	queryResult := s.db.
		Scopes(PaginationScope(options.PaginationOptions)).
		Find(&results)

	<-done

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	paginatedResults := createPaginatedResults(options.PaginationOptions, count, results)
	return paginatedResults, nil
}
