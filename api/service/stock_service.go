package service

import (
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
