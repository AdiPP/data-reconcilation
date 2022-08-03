package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Source defines the properties of a source
type Source struct {
	ID     string
	Amount int
	Desc   string
	Date   time.Time
}

// Proxy defines the properties of a proxy
type Proxy struct {
	ID     string
	Amount int
	Desc   string
	Date   time.Time
}

// Report defines the properties of a report
type Report struct {
	Proxy   Proxy
	Source  Source
	Remarks Remark
}

type Remark struct {
	Discrepancies []string
}

func (r *Report) Array() []string {
	return []string{
		r.Proxy.ID,
		strconv.Itoa(r.Proxy.Amount),
		r.Proxy.Desc,
		r.Proxy.Date.Format("2006-01-02"),
		strings.Join(r.Remarks.Discrepancies, "; "),
	}
}

type CSVProxy struct {
	Headers []string
	Values  [][]string
}

func NewCSVProxy(proxies []Proxy) CSVProxy {
	var CSVProxy CSVProxy

	CSVProxy.Headers = []string{"Amt", "Descr", "Date", "ID"}

	for _, v := range proxies {
		CSVProxy.Values = append(
			CSVProxy.Values,
			[]string{
				strconv.Itoa(v.Amount),
				v.Desc,
				v.Date.Format("2006-01-02"),
				v.ID,
			},
		)
	}

	return CSVProxy
}

type CSVSource struct {
	Headers []string
	Values  [][]string
}

func NewCSVSource(sources []Source) CSVSource {
	var CSVSource CSVSource

	CSVSource.Headers = []string{"Date", "ID", "Amount", "Description"}

	for _, v := range sources {
		CSVSource.Values = append(
			CSVSource.Values,
			[]string{
				v.Date.Format("2006-01-02"),
				v.ID,
				strconv.Itoa(v.Amount),
				v.Desc,
			},
		)
	}

	return CSVSource
}

type CSVReport struct {
	Headers []string
	Values  [][]string
}

func NewCSVReport(reports []Report) CSVReport {
	var CSVReport CSVReport

	CSVReport.Headers = []string{"ID", "Amount", "Description", "Date", "Remarks"}

	for _, v := range reports {
		CSVReport.Values = append(CSVReport.Values, v.Array())
	}

	return CSVReport
}

type SummaryReport struct {
	Headers              [][]string
	DiscrepanciesHeaders []string
	DiscrepanciesValues  [][]string
}

func NewSummaryReport(reports []Report) SummaryReport {
	headers := getSummaryReportHeaders(reports)
	discrepanciesHeaders := []string{"ID", "Amount", "Description", "Date", "Remarks"}
	discrepanciesValues := getSummaryReportDiscrepanciesValues(reports)

	return SummaryReport{
		Headers:              headers,
		DiscrepanciesHeaders: discrepanciesHeaders,
		DiscrepanciesValues:  discrepanciesValues,
	}
}

func getSummaryReportHeaders(reports []Report) [][]string {
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

func getSummaryReportDiscrepanciesValues(reports []Report) [][]string {
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
