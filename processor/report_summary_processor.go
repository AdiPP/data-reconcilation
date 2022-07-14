package processor

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ProcessSummary(sourcepath string, proxypath string, destinationpath string) error {
	results, err := Process(sourcepath, proxypath)

	if err != nil {
		return err
	}

	filename := generateSummaryReportFilename()
	filepath := fmt.Sprintf(destinationpath+"/%s", filename)
	file, err := os.Create(filepath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	headers := getSummaryReportHeaders(results)

	for _, header := range headers {
		_, err = file.WriteString(fmt.Sprintf("%s \n", strings.Join(header, " ")))

		if err != nil {
			return err
		}
	}

	file.WriteString("\n")
	file.WriteString("DISCREPANCIES \n")
	file.WriteString("============== \n")
	file.WriteString("NUMBER | ID | TYPE \n")

	disperencies := getSummaryReportDiscrepancies(results)

	for _, disperency := range disperencies {
		_, err = file.WriteString(disperency + "\n")

		if err != nil {
			return err
		}
	}

	fmt.Printf("Successfully generate %s", filename)

	file.Close()

	return nil
}

func generateSummaryReportFilename() string {
	now := time.Now()

	return fmt.Sprintf(
		"report_summary_%s%s.txt",
		now.Format("20060102"),
		now.Format("150405"),
	)
}

func getSummaryReportHeaders(results []Result) [][]string {
	reportDateRange := []string{"REPORT DATE RANGE:", getSummaryDateRange(results)}
	totalSourceRecordProcessed := 0

	for _, value := range results {
		if value.SourceFound {
			totalSourceRecordProcessed++
		}
	}

	totalSourceRecord := []string{"TOTAL SOURCE RECORD PROCESSED:", strconv.Itoa(totalSourceRecordProcessed)}

	return [][]string{
		reportDateRange,
		totalSourceRecord,
	}
}

func getSummaryDateRange(results []Result) string {
	dates := []time.Time{}

	for _, value := range results {
		dates = append(dates, value.Date)
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

func getSummaryReportDiscrepancies(results []Result) []string {
	discrepancies := []string{}
	i := 1

	for _, value := range results {
		remarks := value.Remarks

		if len(remarks) != 0 {
			for _, remarkValue := range value.Remarks {
				discrepancies = append(
					discrepancies,
					strconv.Itoa(i)+" | "+value.ID+" | "+remarkValue.Type,
				)
				i++
			}
		}
	}

	return discrepancies
}
