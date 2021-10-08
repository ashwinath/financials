package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
)

// SymbolService is the interface to the database for the symbol table
type SymbolService struct {
	db *gorm.DB
}

// NewSymbolService creates a new SessionService
func NewSymbolService(db *gorm.DB) *SymbolService {
	return &SymbolService{
		db: db,
	}
}

// Find finds a session by it's symbol
func (s *SymbolService) Find(symbolString string) (*models.Symbol, error) {
	query := s.db.Where("symbol = ?", symbolString)

	var symbol models.Symbol
	err := query.First(&symbol).Error
	if err != nil {
		return nil, err
	}

	return &symbol, nil
}

// Save saves the session into the database
func (s *SymbolService) Save(symbol *models.Symbol) error {
	return s.db.Save(symbol).Error
}
