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

// SymbolListOptions lists the symbols.
type SymbolListOptions struct {
	PaginationOptions
	SymbolType models.SymbolType
}

func (s *SymbolService) parsesSymbolListOptions(options SymbolListOptions) *gorm.DB {
	return s.db.Where("symbol_type = ?", options.SymbolType)
}

// List lists all trades
func (s *SymbolService) List(options SymbolListOptions) (*PaginatedResults, error) {
	done := make(chan struct{}, 1)

	var count int64
	go func() {
		s.parsesSymbolListOptions(options).Model(&models.Symbol{}).Count(&count)
		done <- struct{}{}
	}()

	var results []models.Symbol
	queryResult := s.parsesSymbolListOptions(options).
		Scopes(PaginationScope(options.PaginationOptions)).
		Find(&results)
	<-done

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	paginatedResults := createPaginatedResults(options.PaginationOptions, count, results)
	return paginatedResults, nil
}

// Save saves the session into the database
func (s *SymbolService) Save(symbol *models.Symbol) error {
	return s.db.Save(symbol).Error
}