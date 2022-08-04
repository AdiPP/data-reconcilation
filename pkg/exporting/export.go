package exporting

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type CSVReportExport struct {
	Headers []string
	Values  [][]string
}

func NewCSVReportExport(reports []Report) CSVReportExport {
	var export CSVReportExport

	export.Headers = []string{"ID", "Amount", "Description", "Date", "Remarks"}

	for _, v := range reports {
		export.Values = append(export.Values, v.Array())
	}

	return export
}

type TextSummaryReportExport struct {
	Headers              [][]string
	DiscrepanciesHeaders []string
	DiscrepanciesValues  [][]string
}

func NewTextSummaryReportExport(reports []Report) TextSummaryReportExport {
	headers := getTextSummaryReportExportHeaders(reports)
	discrepanciesHeaders := []string{"ID", "Amount", "Description", "Date", "Remarks"}
	discrepanciesValues := getTextSummaryReportExportDiscrepanciesValues(reports)

	return TextSummaryReportExport{
		Headers:              headers,
		DiscrepanciesHeaders: discrepanciesHeaders,
		DiscrepanciesValues:  discrepanciesValues,
	}
}

func getTextSummaryReportExportHeaders(reports []Report) [][]string {
	var headers [][]string

	headers = append(
		headers,
		[]string{
			"REPORT DATE RANGE:",
			getSummaryReportDateRange(reports),
		},
	)

	headers = append(
		headers,
		[]string{
			"TOTAL SOURCE RECORDS PROCESSED:",
			getSummaryTotalSourceRecordsProcessed(reports),
		},
	)

	return headers
}

func getTextSummaryReportExportDiscrepanciesValues(reports []Report) [][]string {
	var values [][]string

	for _, v := range reports {
		if len(v.Remarks.Discrepancies) != 0 {
			values = append(values, v.Array())
		}
	}

	return values
}

func getSummaryReportDateRange(reports []Report) string {
	dates := []time.Time{}

	for _, v := range reports {
		dates = append(dates, v.Proxy.Date)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	return fmt.Sprintf(
		"%s - %s",
		dates[1].Format("2006-01-02"),
		dates[len(dates)-1].Format("2006-01-02"),
	)
}

func getSummaryTotalSourceRecordsProcessed(reports []Report) string {
	var total int

	for _, v := range reports {
		if v.Source.ID != "" {
			total++
		}
	}

	return strconv.Itoa(total)
}
