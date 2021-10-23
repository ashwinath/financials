package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashwinath/financials/api/config"
	appcontext "github.com/ashwinath/financials/api/context"
	"github.com/ashwinath/financials/api/controller"
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

	//go appctx.TradeMediator.ProcessTrades()

	router := controller.MakeRouter(appctx)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", c.Server.Port),
		WriteTimeout: c.Server.WriteTimeoutInSeconds,
		ReadTimeout:  c.Server.ReadTimeoutInSeconds,
	}

	go func() {
		log.Printf("Starting server on port: %d", c.Server.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("%s", err.Error())
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server")
	}
}
