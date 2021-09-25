package models

import (
	"time"

	"github.com/lithammer/shortuuid/v3"
	"gorm.io/gorm"
)

// Model contains the base model for every database object
type Model struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate add a short UUID before persisting
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("ID", shortuuid.New())
	return nil
}
