package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
)

// SessionService is the interface to the database for the sessions tabls
type SessionService struct {
	db *gorm.DB
}

// NewSessionService creates a new SessionService
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{
		db: db,
	}
}

// Find finds a session by it's ID
func (s *SessionService) Find(id string) (*models.Session, error) {
	query := s.db.Where("id = ?", id)

	var session models.Session
	err := query.First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Save saves the session into the database
func (s *SessionService) Save(session *models.Session) error {
	return s.db.Debug().Save(session).Error
}

// Delete deletes the session from the database
func (s *SessionService) Delete(session *models.Session) error {
	return s.db.Delete(session).Error
}
