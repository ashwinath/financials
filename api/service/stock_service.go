package service

import (
	"time"

	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StockService is the interface to the database for the sessions tables
type StockService struct {
	db              *gorm.DB
	batchInsertSize int
}

// NewStockService creates a new UserService
func NewStockService(db *gorm.DB, batchInsertSize int) *StockService {
	return &StockService{
		db:              db,
		batchInsertSize: batchInsertSize,
	}
}

// BulkAdd adds multiple transactions at once
func (s *StockService) BulkAdd(transactions []*models.Stock) error {
	return s.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(transactions, s.batchInsertSize).
		Error
}

// Find finds an exchange rate by it's ID
func (s *StockService) Find(symbol string, tradeDate time.Time) (*models.Stock, error) {
	query := s.db.Where("symbol = ?", symbol).Where("trade_date = ?", tradeDate)

	var stonk models.Stock
	err := query.First(&stonk).Error
	if err != nil {
		return nil, err
	}

	return &stonk, nil
}
