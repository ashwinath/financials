package models

import (
	"time"

	"github.com/lithammer/shortuuid/v3"
	"gorm.io/gorm"
)

type Model struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.ID = shortuuid.New()
	return nil
}
