package service

import (
	"time"

	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PortfolioService is the interface to the database for the sessions tables
type PortfolioService struct {
	db              *gorm.DB
	batchInsertSize int
}

// NewPortfolioService creates a new UserService
func NewPortfolioService(db *gorm.DB, batchInsertSize int) *PortfolioService {
	return &PortfolioService{
		db:              db,
		batchInsertSize: batchInsertSize,
	}
}

// BulkAdd adds multiple transactions at once
func (s *PortfolioService) BulkAdd(portfolios []models.Portfolio) error {
	return s.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "trade_date"}, {Name: "symbol"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"principal",
				"nav",
				"quantity",
				"simple_returns",
			}),
		}).
		CreateInBatches(portfolios, s.batchInsertSize).
		Error
}

// List returns the portfolio for a particular user with a time constraint
func (s *PortfolioService) List(userID string, from *time.Time) ([]models.Portfolio, error) {
	var portfolio []models.Portfolio
	err := s.db.
		Where("user_id = ?", userID).
		Where("trade_date >= ?", from).
		Find(&portfolio).
		Error

	if err != nil {
		return nil, err
	}

	return portfolio, err
}
