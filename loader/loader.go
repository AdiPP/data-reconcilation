package loader

import (
	"errors"
	"os"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return file, errors.New(FILE_NOT_FOUND)
	}

	defer file.Close()

	return file, nil
}
