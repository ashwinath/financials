package service

import (
	"testing"

	"github.com/ashwinath/financials/api/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCRUD(t *testing.T) {
	db, err := CreateTestDB()
	assert.Nil(t, err)

	svc := NewUserService(db)

	user := &models.User{
		Username: "ashwinath",
		Password: "myverylongandverysecretpassword",
	}
	err = svc.Save(user)
	assert.Nil(t, err)
	assert.NotNil(t, user.ID)
	defer svc.Delete(user)

	assert.Equal(t, "ashwinath", user.Username)
	assert.Equal(t, "", user.Password)
	assert.NotEqual(t, "", user.PasswordHash)
	assert.NotEqual(t, "myverylongandverysecretpassword", user.PasswordHash)

	found, err := svc.Find(user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, found.ID)
	assert.Equal(t, "", found.Password)
	assert.NotEqual(t, "", found.PasswordHash)
	assert.NotEqual(t, "myverylongandverysecretpassword", found.PasswordHash)

	// try comparing hash
	err = bcrypt.CompareHashAndPassword(
		[]byte(found.PasswordHash),
		[]byte("myverylongandverysecretpassword"),
	)
	assert.Nil(t, err)
}
