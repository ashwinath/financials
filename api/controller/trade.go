package controller

import (
	"net/http"
	"time"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
	"github.com/gorilla/mux"
)

type tradeTransactionController struct {
	controller
}

type bulkTransactionsRequest struct {
	Transactions []*models.Trade `json:"transactions" validate:"required,dive"`
}

func (c *tradeTransactionController) List(w http.ResponseWriter, r *http.Request) {
	session, err := c.getSessionFromCookie(r)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
		return
	}

	options := service.TradeListOptions{}
	if err := c.getParams(r, &options); err != nil {
		badRequest(w, "unable to parse params", err.Error())
		return
	}

	results, err := c.context.TradeMediator.ListTrades(session, options)
	if err != nil {
		internalServiceError(w, "query trades", err.Error())
		return
	}

	ok(w, results)
}

func (c *tradeTransactionController) CreateTransactionInBulk(w http.ResponseWriter, r *http.Request) {
	session, err := c.getSessionFromCookie(r)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
		return
	}

	tx := bulkTransactionsRequest{}
	if err := c.getBody(r, &tx); err != nil {
		badRequest(w, "unable to unmarshal data", "Please check if you are inserting the right values.")
		return
	}

	result := c.context.TradeMediator.CreateTransactionInBulk(session, tx.Transactions)
	if result != nil {
		internalServiceError(w, "error insert", "Error inserting trade transactions.")
		return
	}

	created(w, struct{}{})
}

type listPortfolioParams struct {
	From *time.Time `schema:"from" validate:"required"`
}

type portfolioResult struct {
	Results []models.Portfolio `json:"results"`
}

func (c *tradeTransactionController) ListPortfolio(w http.ResponseWriter, r *http.Request) {
	session, err := c.getSessionFromCookie(r)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
		return
	}

	options := listPortfolioParams{}
	if err := c.getParams(r, &options); err != nil {
		badRequest(w, "unable to parse params", err.Error())
		return
	}

	portfolio, err := c.context.TradeMediator.ListPortfolio(session.UserID, options.From)
	if err != nil {
		internalServiceError(w, "error insert", "Error inserting trade transactions.")
		return
	}

	ok(w, portfolioResult{Results: portfolio})
}

func (c *tradeTransactionController) Delete(w http.ResponseWriter, r *http.Request) {
	session, err := c.getSessionFromCookie(r)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		badRequest(w, "id must be provided", err.Error())
	}

	err = c.context.TradeMediator.Delete(id, session.UserID)
	if err != nil {
		internalServiceError(w, "error insert", "Error deleting trade transaction.")
		return
	}

	ok(w, struct{}{})
}
