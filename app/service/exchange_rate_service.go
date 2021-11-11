package service

import (
	"time"

	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ExchangeRateService is the interface to the database for the sessions tables
type ExchangeRateService struct {
	db              *gorm.DB
	batchInsertSize int
}

// NewExchangeRateService creates a new UserService
func NewExchangeRateService(db *gorm.DB, batchInsertSize int) *ExchangeRateService {
	return &ExchangeRateService{
		db:              db,
		batchInsertSize: batchInsertSize,
	}
}

// BulkAdd adds multiple transactions at once
func (s *ExchangeRateService) BulkAdd(transactions []*models.ExchangeRate) error {
	return s.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(transactions, s.batchInsertSize).
		Error
}

// Find finds an exchange rate by it's ID
func (s *ExchangeRateService) Find(symbol string, tradeDate time.Time) (*models.ExchangeRate, error) {
	query := s.db.Where("symbol = ?", symbol).Where("trade_date = ?", tradeDate)

	var er models.ExchangeRate
	err := query.First(&er).Error
	if err != nil {
		return nil, err
	}

	return &er, nil
}
