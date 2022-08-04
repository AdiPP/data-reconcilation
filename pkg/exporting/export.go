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
	Queries       [][]string
	Summaries     [][]string
	Discrepancies [][]string
}

func NewTextSummaryReportExport(reportOpt ReportOption, reports []Report) TextSummaryReportExport {
	queries := getTextSummaryReportExportQueries(reportOpt, reports)
	summaries := getTextSummaryReportExportSummaries(reports)
	discrepancies := getTextSummaryReportExportDiscrepancies(reports)

	return TextSummaryReportExport{
		Queries:       queries,
		Summaries:     summaries,
		Discrepancies: discrepancies,
	}
}

func getTextSummaryReportExportQueries(reportOpt ReportOption, reports []Report) [][]string {
	var queries [][]string
	var startDate string
	var endDate string

	if !reportOpt.StartDate.IsZero() {
		startDate = reportOpt.StartDate.Format("2006-01-02")
	}

	if !reportOpt.EndDate.IsZero() {
		endDate = reportOpt.EndDate.Format("2006-01-02")
	}

	queries = append(
		queries,
		[]string{
			"Start Date:",
			startDate,
		},
	)

	queries = append(
		queries,
		[]string{
			"End Date:",
			endDate,
		},
	)

	return queries
}

func getTextSummaryReportExportSummaries(reports []Report) [][]string {
	var summaries [][]string

	summaries = append(
		summaries,
		[]string{
			"Proxy Record Date Range:",
			getProxyRecordsDateRange(reports),
		},
	)

	summaries = append(
		summaries,
		[]string{
			"Total Source Records Processed:",
			getTotalSourceRecordsFound(reports),
		},
	)

	discrepancyRecordsCount := 0

	for _, v := range reports {
		if len(v.Remarks) != 0 {
			discrepancyRecordsCount++
		}
	}

	summaries = append(
		summaries,
		[]string{
			"Total Discrepancy Records:",
			strconv.Itoa(discrepancyRecordsCount),
		},
	)

	return summaries
}

func getTextSummaryReportExportDiscrepancies(reports []Report) [][]string {
	var result [][]string

	discrepancies := make(map[string]int)

	for _, v := range reports {
		for _, v := range v.Remarks {

			if val, ok := discrepancies[v.Type]; ok {
				discrepancies[v.Type] = val + 1
			} else {
				discrepancies[v.Type] = 1
			}
		}
	}

	for key, v := range discrepancies {
		result = append(
			result,
			[]string{
				key + ":",
				strconv.Itoa(v),
			},
		)
	}

	return result
}

func getProxyRecordsDateRange(reports []Report) string {
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

func getTotalSourceRecordsFound(reports []Report) string {
	var total int

	for _, v := range reports {
		if v.Source.ID != "" {
			total++
		}
	}

	return strconv.Itoa(total)
}
