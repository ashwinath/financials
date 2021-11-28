package main

import (
	"log"
	"os"

	"github.com/ashwinath/financials/api/config"
	appcontext "github.com/ashwinath/financials/api/context"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Panic(err.Error())
	}

	appctx, err := appcontext.InitContext(c)
	if err != nil {
		log.Panic(err.Error())
	}

	err = appctx.TradeMediator.ProcessTrades()
	if err != nil {
		log.Fatalf("Could not process trades: %s.", err)
		os.Exit(1)
	}
	err = appctx.ExpenseMediator.ProcessExpenses()
	if err != nil {
		log.Fatalf("Could not process expenses: %s.", err)
		os.Exit(1)
	}
}
