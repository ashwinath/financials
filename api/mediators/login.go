package mediator

import (
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
		return nil, err
	}

	return m.createSession(user)
}

// Login creates an account and returns the session id
func (m *LoginMediator) Login(request *models.User) (*models.Session, error) {
	user, err := m.userService.FindByUsername(request.Username)
	if err != nil {
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
