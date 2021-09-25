package mediator

import (
	"errors"
	"testing"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndLogin(t *testing.T) {
	db, err := service.CreateTestDB()
	assert.Nil(t, err)

	sessionService := service.NewSessionService(db)
	userService := service.NewUserService(db)

	loginMediator := NewLoginMediator(userService, sessionService)

	tests := []struct {
		name       string
		createUser *models.User
		loginUser  *models.User
		success    bool
	}{
		{
			name: "success | nominal",
			createUser: &models.User{
				Username: "foo-user",
				Password: "password",
			},
			loginUser: &models.User{
				Username: "foo-user",
				Password: "password",
			},
			success: true,
		},
		{
			name: "failure | wrong password",
			createUser: &models.User{
				Username: "foo-user",
				Password: "password",
			},
			loginUser: &models.User{
				Username: "foo-user",
				Password: "totally-wrong-password",
			},
			success: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := loginMediator.CreateAccount(tt.createUser)
			defer userService.Delete(tt.createUser)
			defer sessionService.Delete(session)

			assert.Nil(t, err)
			assert.NotNil(t, session)
			assert.NotEqual(t, "", session.ID)

			loginSession, err := loginMediator.Login(tt.loginUser)
			if tt.success {
				assert.Nil(t, err)
				assert.NotNil(t, loginSession)
				assert.NotEqual(t, "", loginSession.ID)
			} else {
				assert.True(t, errors.Is(err, ErrorWrongPassword))
			}
		})
	}
}
