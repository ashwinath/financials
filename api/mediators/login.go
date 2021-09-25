package mediator

import (
	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
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
		return nil, err
	}

	session := &models.Session{
		UserID: user.ID,
	}
	if err := m.sessionService.Save(session); err != nil {
		return nil, err
	}

	return session, nil
}
