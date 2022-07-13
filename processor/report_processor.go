package processor

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AdiPP/recon-general/loader"
	"github.com/AdiPP/recon-general/mapper"
)

type Result struct {
	ID          string
	Amount      string
	Description string
	Date        time.Time
	Remarks     []string
}

func Process(sourcepath string, proxypath string) []Result {
	results := []Result{}
	sourceData, err := getSourceData(sourcepath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proxyData, err := getProxyData(proxypath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, proxyValue := range proxyData {
		sourceValue := mapper.SourceValue{}

		for _, value := range sourceData {
			if value.ID == proxyValue.ID {
				sourceValue = value
				break
			}
		}

		results = append(results, Result{
			ID:          proxyValue.ID,
			Amount:      proxyValue.Amount,
			Description: proxyValue.Description,
			Date:        proxyValue.Date,
			Remarks:     getRemarks(proxyValue, sourceValue),
		})
	}

	return results
}

func getSourceData(filepath string) ([]mapper.SourceValue, error) {
	result := []mapper.SourceValue{}
	file, err := loader.NewLoader().Load(filepath)

	if err != nil {
		return result, err
	}

	defer file.Close()

	result, err = mapper.NewSourceMapper().Map(file)

	if err != nil {
		return result, err
	}

	return result, nil
}

func getProxyData(filepath string) ([]mapper.ProxyValue, error) {
	result := []mapper.ProxyValue{}
	file, err := loader.NewLoader().Load(filepath)

	if err != nil {
		return result, err
	}

	defer file.Close()

	result, err = mapper.NewProxyMapper().Map(file)

	if err != nil {
		return result, err
	}

	return result, nil
}

func getRemarks(proxyValue mapper.ProxyValue, sourceValue mapper.SourceValue) []string {
	remarks := make([]string, 0)

	if sourceValue.ID == "" {
		remarks = append(remarks, "Data not found on source")

		return remarks
	}

	if sourceValue.Amount != proxyValue.Amount {
		remarks = append(remarks, "Different amount with source data")
	}

	if sourceValue.Description != proxyValue.Description {
		remarks = append(remarks, "Different description with source data")
	}

	if sourceValue.Date.Format("2006-01-02") != proxyValue.Date.Format("2006-01-02") {
		remarks = append(remarks, "Different date with source data")
	}

	return remarks
}

func ProcessCSV(sourcepath string, proxypath string) {
	results := Process(sourcepath, proxypath)
	filename := fmt.Sprintf(
		"report_%s%s.csv",
		time.Now().Format("20060102"),
		time.Now().Format("150405"),
	)
	filepath := fmt.Sprintf("../result_test/%s", filename)
	file, err := os.Create(filepath)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	values := generateCSVRecord(results)

	for _, value := range values {
		_ = writer.Write(value)
	}

	fmt.Printf("Successfully generate %s", filename)

	writer.Flush()
	file.Close()
}

func generateCSVRecord(results []Result) [][]string {
	headers := []string{"ID", "Amount", "Description", "Date", "Remarks"}
	data := [][]string{headers}

	for _, value := range results {
		data = append(data, []string{
			value.ID,
			value.Amount,
			value.Description,
			value.Date.Format("2006-01-02"),
			strings.Join(value.Remarks, ","),
		})
	}

	return data
}

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
