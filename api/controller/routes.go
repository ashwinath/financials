package controller

import (
	"net/http"

	"github.com/ashwinath/financials/api/context"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type routes struct {
	path    string
	method  string
	handler func(w http.ResponseWriter, _ *http.Request)
}

func makeRoutes(ctx *context.Context) []routes {
	validate := validator.New()
	decoder := schema.NewDecoder()

	base := controller{
		context:   ctx,
		decoder:   decoder,
		validator: validate,
	}

	health := healthController{controller: base}
	login := loginController{controller: base}

	return []routes{
		{"/alive", http.MethodGet, health.Alive},
		{"/ready", http.MethodGet, health.Ready},
		{"/api/v1/users", http.MethodPost, login.CreateUser},
		{"/api/v1/login", http.MethodPost, login.Login},
		{"/api/v1/session", http.MethodGet, login.GetUserFromSession},
	}
}

// MakeRouter makes a multiplexed router
func MakeRouter(ctx *context.Context) *mux.Router {
	r := mux.NewRouter()
	for _, route := range makeRoutes(ctx) {
		r.HandleFunc(route.path, route.handler).Methods(route.method)
	}
	return r
}
