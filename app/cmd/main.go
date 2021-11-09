package main

import (
	"log"

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

	appctx.TradeMediator.ProcessTrades()
}
