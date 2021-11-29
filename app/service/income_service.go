package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IncomeService is the interface to the database for the incomes tables
type IncomeService struct {
	db              *gorm.DB
	batchInsertSize int
}

// NewIncomeService creates a new IncomeService
func NewIncomeService(db *gorm.DB, batchInsertSize int) *IncomeService {
	return &IncomeService{
		db:              db,
		batchInsertSize: batchInsertSize,
	}
}

// TruncateTable truncates the income table
func (s *IncomeService) TruncateTable() error {
	return s.db.Exec("TRUNCATE TABLE incomes;").Error
}

// BulkAdd adds multiple transactions at once
func (s *IncomeService) BulkAdd(income []*models.Income) error {
	return s.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "transaction_date"}, {Name: "type"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"amount",
			}),
		}).
		CreateInBatches(income, s.batchInsertSize).
		Error
}
