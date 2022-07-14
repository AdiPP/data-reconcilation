package mapper

import (
	"encoding/csv"
	"errors"
	"os"
	"time"
)

type SourceValue struct {
	ID          string
	Amount      string
	Description string
	Date        time.Time
}

type SourceMapper struct{}

func NewSourceMapper() *SourceMapper {
	return &SourceMapper{}
}

func (s *SourceMapper) Map(file *os.File) ([]SourceValue, error) {
	lines, err := csv.NewReader(file).ReadAll()

	if err != nil {
		return []SourceValue{}, err
	}

	return s.getResult(lines)
}

// Validate headers to ensure the file are supported
func (s *SourceMapper) validateHeaders(headers []string) error {
	definedHeaders := []string{"Date", "ID", "Amount", "Description"}

	for i, header := range headers {
		if header != definedHeaders[i] {
			return errors.New(FILE_NOT_SUPPORTED)
		}
	}

	return nil
}

func (s *SourceMapper) getResult(lines [][]string) ([]SourceValue, error) {
	result := []SourceValue{}

	for i, line := range lines {
		if i == 0 {
			err := s.validateHeaders(line)

			if err != nil {
				return result, err
			}

			continue
		}

		date, _ := time.Parse("2006-01-02", line[0])

		valueMapper := SourceValue{
			ID:          line[1],
			Amount:      line[2],
			Description: line[3],
			Date:        date,
		}

		result = append(result, valueMapper)
	}

	return result, nil
}
