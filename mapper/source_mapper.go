package mapper

import (
	"encoding/csv"
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

	// Todo-Adi: Add format validationLoader

	result := []SourceValue{}

	for i, line := range lines {
		if i != 0 {
			date, _ := time.Parse(layoutISO, line[0])

			valueMapper := SourceValue{
				ID:          line[1],
				Amount:      line[2],
				Description: line[3],
				Date:        date,
			}

			result = append(result, valueMapper)
		}
	}

	return result, nil
}
