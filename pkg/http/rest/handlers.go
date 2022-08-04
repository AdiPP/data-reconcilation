package rest

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/AdiPP/reconciliation/pkg/exporting"
	"github.com/julienschmidt/httprouter"
)

func Handler(e exporting.Service) http.Handler {
	router := httprouter.New()

	router.GET("/reconcile", reconcile(e))
	router.GET("/reconcile-report", reconcileReport(e))
	router.GET("/reconcile-report-summary", reconcileSummaryReport(e))

	return router
}

func reconcile(e exporting.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		reports, err := e.GetReportData(exporting.ReportOption{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reports)
	}
}

func reconcileReport(e exporting.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		reports, err := e.GetReportData(exporting.ReportOption{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpFile, err := os.CreateTemp(os.TempDir(), "reports-*.csv")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		writer := csv.NewWriter(tmpFile)

		defer writer.Flush()

		err = e.WriteReportFile(writer, reports)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fileBytes, err := ioutil.ReadFile(tmpFile.Name())

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
	}
}

func reconcileSummaryReport(e exporting.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		opt := exporting.ReportOption{}
		reports, err := e.GetReportData(opt)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpFile, err := ioutil.TempFile("", "summary-*.txt")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		err = e.WriteSummaryReportFile(opt, tmpFile, reports)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fileBytes, err := ioutil.ReadFile(tmpFile.Name())

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
	}
}
