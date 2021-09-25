package models

import (
	"time"

	"gorm.io/gorm"
)

const hoursInOneDay = 24

// Session contains the login session of a user.
type Session struct {
	Model
	UserID string
	Expiry *time.Time
}

// BeforeCreate adds the exprity date for the session
func (m *Session) BeforeCreate(tx *gorm.DB) error {
	// Need to override the models Before Create as they don't stack
	err := m.Model.BeforeCreate(tx)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	expiry := time.Now().In(loc).Add(time.Hour * hoursInOneDay)
	m.Expiry = &expiry
	return nil
}
