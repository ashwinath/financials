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
	c := healthController{
		controller: controller{
			context:   ctx,
			decoder:   decoder,
			validator: validate,
		},
	}
	return []routes{
		{"/alive", http.MethodGet, c.Alive},
		{"/ready", http.MethodGet, c.Ready},
	}
}

func MakeRouter(ctx *context.Context) *mux.Router {
	r := mux.NewRouter()
	for _, route := range makeRoutes(ctx) {
		r.HandleFunc(route.path, route.handler).Methods(route.method)
	}
	return r
}
