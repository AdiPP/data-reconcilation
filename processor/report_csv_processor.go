package processor

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

func ProcessCSV(sourcepath string, proxypath string, destinationpath string) error {
	results, err := Process(sourcepath, proxypath)

	if err != nil {
		return err
	}

	filename := generateCSVReportFilename()
	filepath := fmt.Sprintf(destinationpath+"/%s", filename)
	file, err := os.Create(filepath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	values := getCSVValue(results)

	for _, value := range values {
		_ = writer.Write(value)
	}

	writer.Flush()
	file.Close()

	return nil
}

func generateCSVReportFilename() string {
	now := time.Now()

	return fmt.Sprintf(
		"report_%s%s.csv",
		now.Format("20060102"),
		now.Format("150405"),
	)
}

func getCSVValue(results []Result) [][]string {
	headers := []string{"ID", "Amount", "Description", "Date", "Remarks"}
	data := [][]string{headers}

	for _, value := range results {
		data = append(data, []string{
			value.ID,
			value.Amount,
			value.Description,
			value.Date.Format("2006-01-02"),
			convertRemarksToCSVValue(value.Remarks),
		})
	}

	return data
}

func convertRemarksToCSVValue(remarks []Remark) string {
	messages := make([]string, 0)

	for _, value := range remarks {
		messages = append(messages, value.Message)
	}

	return strings.Join(messages, ",")
}
