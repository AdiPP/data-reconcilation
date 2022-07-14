package processor

import (
	"fmt"
	"time"

	"github.com/AdiPP/recon-general/loader"
	"github.com/AdiPP/recon-general/mapper"
)

type Result struct {
	ID          string
	Amount      string
	Description string
	Date        time.Time
	SourceFound bool
	Remarks     []Remark
}

type Remark struct {
	Type    string
	Message string
}

func Process(sourcepath string, proxypath string) ([]Result, error) {
	results := []Result{}
	sourceData, err := getSourceData(sourcepath)

	if err != nil {
		return results, err
	}

	proxyData, err := getProxyData(proxypath)

	if err != nil {
		return results, err
	}

	for _, proxyValue := range proxyData {
		sourceValue := mapper.SourceValue{}
		sourceFound := false

		for _, value := range sourceData {
			if value.ID == proxyValue.ID {
				sourceValue = value
				sourceFound = true
				break
			}
		}

		results = append(results, Result{
			ID:          proxyValue.ID,
			Amount:      proxyValue.Amount,
			Description: proxyValue.Description,
			Date:        proxyValue.Date,
			SourceFound: sourceFound,
			Remarks:     getRemarks(sourceFound, proxyValue, sourceValue),
		})
	}

	return results, nil
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

func getRemarks(sourceFound bool, proxyValue mapper.ProxyValue, sourceValue mapper.SourceValue) []Remark {
	remarks := make([]Remark, 0)

	if !sourceFound {
		remarks = append(remarks, Remark{
			Type:    "DATA_NOT_FOUND",
			Message: DATA_NOT_FOUND,
		})

		return remarks
	}

	if sourceValue.Amount != proxyValue.Amount {
		remarks = append(remarks, Remark{
			Type:    "DIFFERENT_AMOUNT",
			Message: fmt.Sprintf(DIFFERENT_AMOUNT, proxyValue.Amount, sourceValue.Amount),
		})
	}

	if sourceValue.Description != proxyValue.Description {
		remarks = append(remarks, Remark{
			Type:    "DIFFERENT_DESCRIPTION",
			Message: fmt.Sprintf(DIFFERENT_DESCRIPTION, proxyValue.Description, sourceValue.Description),
		})
	}

	sourceDate := sourceValue.Date.Format("2006-01-02")
	proxyDate := proxyValue.Date.Format("2006-01-02")

	if sourceDate != proxyDate {
		remarks = append(remarks, Remark{
			Type:    "DIFFERENT_DATE",
			Message: fmt.Sprintf(DIFFERENT_DATE, proxyDate, sourceDate),
		})
	}

	return remarks
}
