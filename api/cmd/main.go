package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ashwinath/financials/api/config"
	"github.com/ashwinath/financials/api/context"
	"github.com/ashwinath/financials/api/controller"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Panic(err.Error())
	}

	ctx, err := context.InitContext(c)
	if err != nil {
		log.Panic(err.Error())
	}

	router := controller.MakeRouter(ctx)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", c.Server.Port),
		WriteTimeout: c.Server.WriteTimeoutInSeconds,
		ReadTimeout:  c.Server.ReadTimeoutInSeconds,
	}

	log.Printf("Starting server on port: %d", c.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
