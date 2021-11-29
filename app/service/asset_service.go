package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AssetService is the interface to the database for the sessions tables
type AssetService struct {
	db              *gorm.DB
	batchInsertSize int
}

// NewAssetService creates a new AssetService
func NewAssetService(db *gorm.DB, batchInsertSize int) *AssetService {
	return &AssetService{
		db:              db,
		batchInsertSize: batchInsertSize,
	}
}

// TruncateTable truncates the assets table
func (s *AssetService) TruncateTable() error {
	return s.db.Exec("TRUNCATE TABLE assets;").Error
}

// BulkAdd adds multiple transactions at once
func (s *AssetService) BulkAdd(expenses []*models.Asset) error {
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
