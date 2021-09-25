package models

import (
	"time"

	"gorm.io/gorm"
)

const HoursInOneDay = 24

type Session struct {
	Model
	UserID string
	Expiry *time.Time
}

func (m *Session) BeforeCreate(tx *gorm.DB) error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	expiry := time.Now().In(loc).Add(time.Hour * HoursInOneDay)
	m.Expiry = &expiry
	return nil
}
