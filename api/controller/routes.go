package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/ashwinath/financials/api/context"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const indexPath = "index.html"

type routes struct {
	path    string
	method  string
	handler func(http.ResponseWriter, *http.Request)
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
	trades := tradeTransactionController{controller: base}

	return []routes{
		{"/alive", http.MethodGet, health.Alive},
		{"/ready", http.MethodGet, health.Ready},
		{"/api/v1/users", http.MethodPost, login.CreateUser},
		{"/api/v1/login", http.MethodPost, login.Login},
		{"/api/v1/logout", http.MethodPost, login.Logout},
		{"/api/v1/session", http.MethodGet, login.GetUserFromSession},
		{"/api/v1/trades", http.MethodGet, trades.List},
		{"/api/v1/trades", http.MethodPost, trades.CreateTransactionInBulk},
		{"/api/v1/trades/portfolio", http.MethodGet, trades.ListPortfolio},
	}
}

func makeFrontendHandler(ctx *context.Context) func(http.ResponseWriter, *http.Request) {
	frontendPath := ctx.Config.Server.ReactFilePath
	frontendHandler := http.FileServer(http.Dir(frontendPath))
	return func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		path = filepath.Join(frontendPath, path)
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			// file does not exist, serve index.html
			http.ServeFile(w, r, filepath.Join(frontendPath, indexPath))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		frontendHandler.ServeHTTP(w, r)
	}
}

// MakeRouter makes a multiplexed router
func MakeRouter(ctx *context.Context) *mux.Router {
	r := mux.NewRouter()
	for _, route := range makeRoutes(ctx) {
		r.HandleFunc(route.path, route.handler).Methods(route.method)
	}
	r.PathPrefix("/").HandlerFunc(makeFrontendHandler(ctx))

	return r
}
