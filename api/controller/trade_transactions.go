package controller

import (
	"net/http"

	"github.com/ashwinath/financials/api/models"
	"github.com/ashwinath/financials/api/service"
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
