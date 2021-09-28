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

// Find finds a session by it's ID
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
func (s *TradeService) BulkAdd(transactions []models.Trade) error {
	return s.db.CreateInBatches(transactions, s.batchInsertSize).Error
}

// Delete deletes a transaction
func (s *TradeService) Delete(transaction *models.Trade) error {
	return s.db.Delete(transaction).Error
}
