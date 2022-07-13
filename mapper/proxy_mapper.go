package mapper

import (
	"encoding/csv"
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

func (m *ProxyMapper) Map(file *os.File) ([]ProxyValue, error) {
	lines, err := csv.NewReader(file).ReadAll()

	if err != nil {
		return []ProxyValue{}, err
	}

	// Todo-Adi: Add format validation

	result := []ProxyValue{}

	for i, line := range lines {
		if i != 0 {
			date, _ := time.Parse(layoutISO, line[2])

			valueMapper := ProxyValue{
				ID:          line[3],
				Amount:      line[0],
				Description: line[1],
				Date:        date,
			}

			result = append(result, valueMapper)
		}
	}

	return result, nil
}
