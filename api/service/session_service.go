package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{
		db: db,
	}
}

func (s *SessionService) Find(id string) (*models.Session, error) {
	query := s.db.Where("id = ?", id)

	var session models.Session
	err := query.First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionService) Save(session *models.Session) error {
	return s.db.Save(session).Error
}

func (s *SessionService) Delete(session *models.Session) error {
	return s.db.Delete(session).Error
}
