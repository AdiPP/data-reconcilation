package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AdiPP/reconciliation/pkg/exporting"
	"github.com/AdiPP/reconciliation/pkg/http/rest"
	"github.com/AdiPP/reconciliation/pkg/storage/memory"
)

const (
	DataDirPath = "../../resources/"
)

func main() {
	var exporter exporting.Service

	s, err := memory.NewStorageFromFile(
		DataDirPath+"proxy.csv",
		DataDirPath+"source.csv",
	)

	if err != nil {
		panic(err.Error())
	}

	exporter = exporting.NewService(s)

	router := rest.Handler(exporter)

	fmt.Println("The beer server is on tap now: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
