package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	FileDataLocation = "./data/"
)

type StorageFile struct {
	proxypath  string
	sourcepath string
	proxys     []Proxy
	sources    []Source
}

func NewStorageFile(proxyPath string, sourcePath string) (*StorageFile, error) {
	var err error

	stg := new(StorageFile)
	stg.proxypath = proxyPath
	stg.sourcepath = sourcePath

	prxyRecords, err := getRecords(stg.proxypath)

	if err != nil {
		return nil, err
	}

	prxys, err := mapProxyRecords(prxyRecords)

	if err != nil {
		return nil, err
	}

	stg.proxys = prxys

	srcRecords, err := getRecords(stg.sourcepath)

	if err != nil {
		return nil, err
	}

	srcs, err := mapSourceRecords(srcRecords)

	if err != nil {
		return nil, err
	}

	stg.sources = srcs

	return stg, nil
}

func (s *StorageFile) FindProxyByID(ID string) (Proxy, error) {
	var proxy Proxy

	for _, v := range s.proxys {
		if v.ID == ID {
			return v, nil
		}
	}

	return proxy, nil
}

func (s *StorageFile) FindAllSources() ([]Source, error) {
	return s.sources, nil
}

func (s *StorageFile) FindSourceByID(ID string) (Source, error) {
	var source Source

	for _, v := range s.sources {
		if v.ID == ID {
			return v, nil
		}
	}

	return source, nil
}

func (s *StorageFile) FindAllProxies() ([]Proxy, error) {
	return s.proxys, nil
}

func getRecords(path string) ([][]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()

	if err != nil {
		return nil, err
	}

	return records, nil
}

func mapProxyRecords(records [][]string) ([]Proxy, error) {
	var prxs []Proxy

	for i, v := range records {
		// Skip csv file header
		if i == 0 {
			err := validateFileHeader(v, []string{"Amt", "Descr", "Date", "ID"})

			if err != nil {
				return nil, err
			}

			continue
		}

		amount, err := strconv.Atoi(v[0])

		if err != nil {
			return nil, err
		}

		date, err := time.Parse("2006-01-02", v[2])

		if err != nil {
			return nil, err
		}

		prxs = append(prxs, Proxy{
			ID:     v[3],
			Amount: amount,
			Desc:   v[1],
			Date:   date,
		})
	}

	return prxs, nil
}

func mapSourceRecords(records [][]string) ([]Source, error) {
	var srcs []Source

	for i, v := range records {
		// Skip csv file header
		if i == 0 {
			err := validateFileHeader(v, []string{"Date", "ID", "Amount", "Description"})

			if err != nil {
				return nil, err
			}

			continue
		}

		amount, err := strconv.Atoi(v[2])

		if err != nil {
			return nil, err
		}

		date, err := time.Parse("2006-01-02", v[0])

		if err != nil {
			return nil, err
		}

		srcs = append(srcs, Source{
			ID:     v[1],
			Amount: amount,
			Desc:   v[3],
			Date:   date,
		})
	}

	return srcs, nil
}

func validateFileHeader(fileHeaders []string, definedHeaders []string) error {
	var err error

	for i, v := range fileHeaders {
		if v != definedHeaders[i] {
			return errors.New("unsupported file header")
		}
	}

	return err
}
