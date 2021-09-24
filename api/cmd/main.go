package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ashwinath/financials/api/config"
	"github.com/gorilla/mux"
)

func main() {
	configuration, err := config.Load()
	if err != nil {
		log.Panic(err.Error())
	}

	r := mux.NewRouter()
	//r.HandleFunc("/foo", foo.DoSomething).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", configuration.Server.Port),
		WriteTimeout: configuration.Server.WriteTimeoutInSeconds,
		ReadTimeout:  configuration.Server.ReadTimeoutInSeconds,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
