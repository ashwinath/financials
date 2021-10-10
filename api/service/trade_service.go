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

// Find finds a trade by it's ID
func (s *TradeService) Find(id string) (*models.Trade, error) {
	query := s.db.Where("id = ?", id)

	var trade models.Trade
	err := query.First(&trade).Error
	if err != nil {
		return nil, err
	}

	return &trade, nil
}

// BulkAdd adds multiple transactions at once
func (s *TradeService) BulkAdd(transactions []*models.Trade) error {
	return s.db.CreateInBatches(transactions, s.batchInsertSize).Error
}

// Delete deletes a transaction
func (s *TradeService) Delete(transaction *models.Trade) error {
	return s.db.Delete(transaction).Error
}

// TradeListOptions lists the trades from a user.
type TradeListOptions struct {
	PaginationOptions
	UserID *string `schema:"-"`
}

func (s *TradeService) parseTradeListOptions(options TradeListOptions) *gorm.DB {
	// This is required as method chaining has some problems with concurrency.
	// See: https://gorm.io/docs/method_chaining.html#Method-Chain-Safety-Goroutine-Safety
	query := s.db

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}

	return query
}

// List lists all trades
func (s *TradeService) List(options TradeListOptions) (*PaginatedResults, error) {
	done := make(chan struct{}, 1)

	var count int64
	go func() {
		s.parseTradeListOptions(options).Model(&models.Trade{}).Count(&count)
		done <- struct{}{}
	}()

	var results []models.Trade
	queryResult := s.parseTradeListOptions(options).
		Scopes(PaginationScope(options.PaginationOptions)).
		Find(&results)

	<-done

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	paginatedResults := createPaginatedResults(options.PaginationOptions, count, results)
	return paginatedResults, nil
}
