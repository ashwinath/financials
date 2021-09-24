package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ashwinath/financials/api/config"
	"github.com/ashwinath/financials/api/context"
	"github.com/gorilla/mux"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Panic(err.Error())
	}

	_, err = context.InitContext(c)
	if err != nil {
		log.Panic(err.Error())
	}

	r := mux.NewRouter()
	//r.HandleFunc("/foo", foo.DoSomething).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", c.Server.Port),
		WriteTimeout: c.Server.WriteTimeoutInSeconds,
		ReadTimeout:  c.Server.ReadTimeoutInSeconds,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
