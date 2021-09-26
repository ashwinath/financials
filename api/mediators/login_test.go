package mediator

import (
	"errors"
	"testing"
	"time"

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
		errorType  error
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
			success:   true,
			errorType: nil,
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
			success:   false,
			errorType: ErrorWrongPassword,
		},
		{
			name: "failure | no such user",
			createUser: &models.User{
				Username: "foo-user",
				Password: "password",
			},
			loginUser: &models.User{
				Username: "foo-user-hello",
				Password: "totally-wrong-password",
			},
			success:   false,
			errorType: ErrorNoSuchUser,
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
				assert.True(t, errors.Is(err, tt.errorType))
			}
		})
	}
}

func TestDuplicateUser(t *testing.T) {
	db, err := service.CreateTestDB()
	assert.Nil(t, err)

	sessionService := service.NewSessionService(db)
	userService := service.NewUserService(db)

	loginMediator := NewLoginMediator(userService, sessionService)
	t.Run("failure | duplicate user", func(t *testing.T) {
		user := &models.User{
			Username: "duplicate",
			Password: "helloworld",
		}
		session, err := loginMediator.CreateAccount(user)
		defer userService.Delete(user)
		defer sessionService.Delete(session)
		assert.Nil(t, err)
		assert.NotNil(t, session)
		assert.NotEqual(t, "", session.ID)

		// Duplicate here
		duplicateUser := &models.User{
			Username: "duplicate",
			Password: "helloworld",
		}
		newSession, err := loginMediator.CreateAccount(duplicateUser)
		assert.Nil(t, newSession)
		assert.Equal(t, err, ErrorDuplicateUser)
	})
}

func TestGetUserFromSession(t *testing.T) {
	db, err := service.CreateTestDB()
	assert.Nil(t, err)

	sessionService := service.NewSessionService(db)
	userService := service.NewUserService(db)
	user := &models.User{
		Username: "duplicate",
		Password: "helloworld",
	}
	loginMediator := NewLoginMediator(userService, sessionService)
	session, err := loginMediator.CreateAccount(user)
	defer userService.Delete(user)
	defer sessionService.Delete(session)

	tests := []struct {
		name      string
		sessionID string
		success   bool
	}{
		{
			name:      "success | nominal",
			sessionID: session.ID,
			success:   true,
		},
		{
			name:      "failure | no such session",
			sessionID: "fake-session",
			success:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := loginMediator.GetUserFromSession(tt.sessionID)
			if tt.success {
				assert.Nil(t, err)
				assert.Equal(t, user.Username, found.Username)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, found)
			}
		})
	}

	t.Run("failure | expired session", func(t *testing.T) {
		loc, err := time.LoadLocation("Asia/Singapore")
		assert.Nil(t, err)
		expiry := time.Now().In(loc).Add(time.Hour * -1)
		session.Expiry = &expiry
		sessionService.Save(session)
		found, err := loginMediator.GetUserFromSession(session.ID)
		assert.NotNil(t, err)
		assert.Nil(t, found)
	})
}

func TestLogUserOut(t *testing.T) {
	db, err := service.CreateTestDB()
	assert.Nil(t, err)

	sessionService := service.NewSessionService(db)
	userService := service.NewUserService(db)

	loginMediator := NewLoginMediator(userService, sessionService)

	user := &models.User{
		Username: "duplicate",
		Password: "helloworld",
	}
	session, err := loginMediator.CreateAccount(user)
	defer userService.Delete(user)

	err = loginMediator.Logout(session.ID)
	assert.Nil(t, err)

	found, err := sessionService.Find(session.ID)
	assert.Nil(t, found)
	assert.NotNil(t, err)
}
