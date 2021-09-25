package service

import (
	"github.com/ashwinath/financials/api/models"
	"gorm.io/gorm"
)

// UserService is the interface to the database for the sessions tabls
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// Find finds a user by it's ID
func (s *UserService) Find(id string) (*models.User, error) {
	query := s.db.Where("id = ?", id)

	var user models.User
	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByUsername finds a user by username
func (s *UserService) FindByUsername(username string) (*models.User, error) {
	query := s.db.Where("username = ?", username)

	var user models.User
	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Save saves the session into the database
func (s *UserService) Save(user *models.User) error {
	return s.db.Save(user).Error
}

// Delete deletes the session from the database
func (s *UserService) Delete(user *models.User) error {
	return s.db.Delete(user).Error
}
