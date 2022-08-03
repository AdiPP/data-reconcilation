package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnsupportedHeaderProxyFile(t *testing.T) {
	_, err := NewStorageFile(
		FileDataLocation+"proxy_invalid_header.csv",
		FileDataLocation+"source.csv",
	)

	assert.NotNil(t, err)
}

func TestUnsupportedHeaderSourceFile(t *testing.T) {
	_, err := NewStorageFile(
		FileDataLocation+"proxy.csv",
		FileDataLocation+"source_invalid_header.csv",
	)

	assert.NotNil(t, err)
}

func TestFindProxyByID(t *testing.T) {
	storage, err := NewStorageFile(
		FileDataLocation+"proxy.csv",
		FileDataLocation+"source.csv",
	)

	assert.Nil(t, err)

	time, err := time.Parse("2006-01-02", "2021-06-30")

	assert.Nil(t, err)

	sampleProxy := Proxy{
		ID:     "zoUr",
		Amount: 24,
		Desc:   "A",
		Date:   time,
	}

	proxy, err := storage.FindProxyByID(sampleProxy.ID)

	assert.Nil(t, err)
	assert.NotNil(t, proxy.ID)
	assert.Equal(t, sampleProxy.ID, proxy.ID)
	assert.Equal(t, sampleProxy.Amount, proxy.Amount)
	assert.Equal(t, sampleProxy.Desc, proxy.Desc)
	assert.Equal(t, sampleProxy.Date, proxy.Date)
}

func TestFindSourceByID(t *testing.T) {
	storage, err := NewStorageFile(
		FileDataLocation+"proxy.csv",
		FileDataLocation+"source.csv",
	)

	assert.Nil(t, err)

	time, err := time.Parse("2006-01-02", "2021-06-30")

	assert.Nil(t, err)

	sampleSource := Proxy{
		ID:     "zoUr",
		Amount: 24,
		Desc:   "A",
		Date:   time,
	}

	proxy, err := storage.FindProxyByID(sampleSource.ID)

	assert.Nil(t, err)
	assert.NotNil(t, proxy.ID)
	assert.Equal(t, sampleSource.ID, proxy.ID)
	assert.Equal(t, sampleSource.Amount, proxy.Amount)
	assert.Equal(t, sampleSource.Desc, proxy.Desc)
	assert.Equal(t, sampleSource.Date, proxy.Date)
}
