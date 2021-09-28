package controller

import (
	"net/http"

	"github.com/ashwinath/financials/api/models"
)

type tradeTransactionController struct {
	controller
}

type bulkTransactionsRequest struct {
	Transactions []models.Trade `json:"transactions" validate:"required"`
}

func (c *tradeTransactionController) CreateTransactionInBulk(w http.ResponseWriter, r *http.Request) {
	session, err := c.getSessionFromCookie(r)
	if err != nil {
		badRequest(w, "session not found", "not a valid session.")
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
