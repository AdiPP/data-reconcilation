package mapper

import (
	"encoding/csv"
	"errors"
	"os"
	"time"
)

type ProxyValue struct {
	ID          string
	Amount      string
	Description string
	Date        time.Time
}

type ProxyMapper struct{}

func NewProxyMapper() *ProxyMapper {
	return &ProxyMapper{}
}

func (p *ProxyMapper) Map(file *os.File) ([]ProxyValue, error) {
	lines, err := csv.NewReader(file).ReadAll()

	if err != nil {
		return []ProxyValue{}, err
	}

	return p.getResult(lines)
}

// Validate headers to ensure the file are supported
func (p *ProxyMapper) validateHeaders(headers []string) error {
	definedHeaders := []string{"Amt", "Descr", "Date", "ID"}

	for i, header := range headers {
		if header != definedHeaders[i] {
			return errors.New(FILE_NOT_SUPPORTED)
		}
	}

	return nil
}

func (p *ProxyMapper) getResult(lines [][]string) ([]ProxyValue, error) {
	result := []ProxyValue{}

	for i, line := range lines {
		if i == 0 {
			err := p.validateHeaders(line)

			if err != nil {
				return []ProxyValue{}, err
			}

			continue
		}

		date, err := time.Parse("2006-01-02", line[2])

		if err != nil {
			return []ProxyValue{}, err
		}

		value := ProxyValue{
			ID:          line[3],
			Amount:      line[0],
			Description: line[1],
			Date:        date,
		}

		result = append(result, value)
	}

	return result, nil
}
