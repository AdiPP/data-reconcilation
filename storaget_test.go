package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatingFileStorage(t *testing.T) {
	s, err := NewStorageFile(
		FileDataLocation+"proxy.csv",
		FileDataLocation+"source.csv",
	)

	assert.Nil(t, err)
	assert.IsType(t, &StorageFile{}, s)
}
