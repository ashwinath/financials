package service

import (
	"testing"

	"github.com/ashwinath/financials/api/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSessionCRUD(t *testing.T) {
	t.Run("success | nominal", func(t *testing.T) {
		WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
			svc := NewSessionService(db)

			session := &models.Session{
				UserID: "hello-world",
			}
			err := svc.Save(session)
			assert.Nil(t, err)
			assert.NotNil(t, session.ID)

			assert.Equal(t, "hello-world", session.UserID)
			assert.NotNil(t, session.Expiry)

			found, err := svc.Find(session.ID)
			assert.Nil(t, err)
			assert.NotNil(t, found.ID)
			assert.Equal(t, "hello-world", found.UserID)
			assert.NotNil(t, found.Expiry)
		})
	})
}
