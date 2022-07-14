package processor

import (
	"fmt"
	"os"
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
