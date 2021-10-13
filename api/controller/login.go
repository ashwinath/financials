package controller

import (
	"errors"
	"net/http"

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
		if errors.Is(err, mediator.ErrorDuplicateUser) {
			badRequest(w, "username taken", "Username has been taken.")
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

type loginResponse struct {
	Username string `json:"username"`
}

func (c *loginController) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := c.getBody(r, user); err != nil {
		badRequest(w, "failed to parse request", err.Error())
		return
	}

	session, err := c.context.LoginMediator.Login(user)
	if err != nil {
		if errors.Is(err, mediator.ErrorNoSuchUser) {
			badRequest(w, "no such user", "User does not exist.")
			return
		}
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

	ok(w, loginResponse{Username: user.Username})
}

func (c *loginController) GetUserFromSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieSessionName)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
		return
	}

	user, err := c.context.LoginMediator.GetUserFromSession(cookie.Value)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
		return
	}

	ok(w, user)
}

func (c *loginController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieSessionName)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
		return
	}

	err = c.context.LoginMediator.Logout(cookie.Value)
	if err != nil {
		badRequest(w, "could not delete session", "Session not found, can't log you out..")
		return
	}

	// Blankco out cookie
	http.SetCookie(w, &http.Cookie{
		Name:  CookieSessionName,
		Value: "",
	})

	ok(w, struct{}{})
}
