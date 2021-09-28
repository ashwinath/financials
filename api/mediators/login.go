package mediator

import (
	"strings"
	"time"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"golang.org/x/crypto/bcrypt"
)

// LoginMediator handles everything regarding logging in and authentication
type LoginMediator struct {
	userService    *service.UserService
	sessionService *service.SessionService
}

// NewLoginMediator creates a new LoginMediator
func NewLoginMediator(
	userService *service.UserService,
	sessionService *service.SessionService,
) *LoginMediator {
	return &LoginMediator{
		userService:    userService,
		sessionService: sessionService,
	}
}

// CreateAccount creates an account and returns the session id
func (m *LoginMediator) CreateAccount(user *models.User) (*models.Session, error) {
	if err := m.userService.Save(user); err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return nil, ErrorDuplicateUser
		}
		return nil, err
	}

	return m.createSession(user)
}

// Login creates an account and returns the session id
func (m *LoginMediator) Login(request *models.User) (*models.Session, error) {
	user, err := m.userService.FindByUsername(request.Username)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, ErrorNoSuchUser
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	)

	if err != nil {
		return nil, ErrorWrongPassword
	}
	return m.createSession(user)
}

func (m *LoginMediator) createSession(user *models.User) (*models.Session, error) {
	session := &models.Session{
		UserID: user.ID,
	}

	if err := m.sessionService.Save(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession gets the userID from the sessionID
func (m *LoginMediator) GetSession(sessionID string) (*models.Session, error) {
	session, err := m.sessionService.Find(sessionID)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().In(loc)
	if session.Expiry.Before(currentTime) {
		return nil, ErrorExpiredSession
	}

	return session, nil
}

// GetUserFromSession checks if the session has expired and returns the user.
func (m *LoginMediator) GetUserFromSession(sessionID string) (*models.User, error) {
	session, err := m.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	user, err := m.userService.Find(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Logout deletes a session for a user.
func (m *LoginMediator) Logout(sessionID string) error {
	session, err := m.sessionService.Find(sessionID)
	if err != nil {
		return err
	}

	return m.sessionService.Delete(session)
}
