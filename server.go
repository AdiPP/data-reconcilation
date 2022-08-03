package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ServerAddr defines the http host and port of the reconcilation server
const ServerAddr = "localhost:8080"

func ExecuteServer() {
	var err error

	PopulateProxies()
	PopulateSources()

	db, err = NewStorageFile(
		"./data/proxy.csv",
		"./data/source.csv",
	)

	if err != nil {
		log.Fatal(err)
	}

	router = httprouter.New()

	router.GET("/reconcile", ReconcileHandler)
	router.GET("/reconcile-csv", ReconcileCSVHandler)
	router.GET("/reconcile-summary", ReconcileSummaryHandler)

	fmt.Println("The reconcilation server is on tap at http://localhost:8080.")
	log.Fatal(http.ListenAndServe(ServerAddr, router))
}
