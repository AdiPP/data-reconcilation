package processor

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

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
