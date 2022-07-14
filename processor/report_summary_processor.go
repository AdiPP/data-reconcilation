package processor

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ProcessSummary(sourcepath string, proxypath string) {
	results := Process(sourcepath, proxypath)
	filename := fmt.Sprintf(
		"report_summary_%s%s.txt",
		time.Now().Format("20060102"),
		time.Now().Format("150405"),
	)
	filepath := fmt.Sprintf("../result_test/%s", filename)
	file, err := os.Create(filepath)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	values := generateSummaryRecord(results)

	for _, value := range values {
		_, err = file.WriteString(fmt.Sprintf("%s \n", strings.Join(value, ": ")))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully generate %s", filename)

	file.Close()
}

func generateSummaryRecord(results []Result) [][]string {
	dateRange := []string{"Date Range", getSummaryDateRange(results)}
	totalRecord := []string{"Total Record", strconv.Itoa(len(results))}
	// disperencies := []string{}

	// for _, value := range results {
	// 	disperencies = append(disperencies, strings.Join(value.Remarks, ","))
	// }

	return [][]string{
		dateRange,
		totalRecord,
		// disperencies,
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
