package controller

import (
	"errors"
	"net/http"
	"strings"

	mediator "github.com/ashwinath/financials/api/mediators"
	"github.com/ashwinath/financials/api/models"
)

type loginController struct {
	controller
}

func (c *loginController) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := c.getBody(r, user); err != nil {
		badRequest(w, "failed to parse request", err.Error())
		return
	}

	session, err := c.context.LoginMediator.CreateAccount(user)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			badRequest(w, "username taken", "Username has been taken, please pick another username.")
			return
		}

		internalServiceError(w, "something went wrong when creating an account", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    CookieSessionName,
		Value:   session.ID,
		Expires: *session.Expiry,
	})

	created(w, user)
}

func (c *loginController) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := c.getBody(r, user); err != nil {
		badRequest(w, "failed to parse request", err.Error())
		return
	}

	session, err := c.context.LoginMediator.Login(user)
	if err != nil {
		if errors.Is(err, mediator.ErrorWrongPassword) {
			badRequest(w, "wrong password", "Password did not match.")
			return
		}

		internalServiceError(w, "something went wrong when trying to login", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    CookieSessionName,
		Value:   session.ID,
		Expires: *session.Expiry,
	})

	ok(w, struct{}{})
}
