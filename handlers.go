package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ReconcileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	reports, err := Reconcile()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func ReconcileCSVHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fileBytes, err := ReconcileCSVFileByte()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func ReconcileSummaryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fileBytes, err := ReconcileSummaryFileByte()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
