package loader

import (
	"os"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return file, err
	}

	return file, nil
}
