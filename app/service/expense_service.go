package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ExpenseService is the interface to the database for the sessions tables
type ExpenseService struct {
	db              *gorm.DB
	batchInsertSize int
}

// NewExpenseService creates a new ExpenseService
func NewExpenseService(db *gorm.DB, batchInsertSize int) *ExpenseService {
	return &ExpenseService{
		db:              db,
		batchInsertSize: batchInsertSize,
	}
}

// TruncateTable truncates the expenses table
func (s *ExpenseService) TruncateTable() error {
	return s.db.Exec("TRUNCATE TABLE expenses;").Error
}

// BulkAdd adds multiple transactions at once
func (s *ExpenseService) BulkAdd(expenses []*models.Expense) error {
	return s.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "transaction_date"}, {Name: "type"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"amount",
			}),
		}).
		CreateInBatches(expenses, s.batchInsertSize).
		Error
}
