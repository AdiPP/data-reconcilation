package main

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"time"
)

func PopulateProxies() {
	defaultProxies := []Proxy{
		{
			ID:     "zoUr",
			Amount: 24,
			Desc:   "A",
			Date:   getDate("2021-06-30"),
		},
		{
			ID:     "zoXq",
			Amount: 11,
			Desc:   "B",
			Date:   getDate("2021-06-30"),
		},
		{
			ID:     "zoap",
			Amount: 69,
			Desc:   "C",
			Date:   getDate("2021-07-01"),
		},
		{
			ID:     "zodo",
			Amount: 30,
			Desc:   "D",
			Date:   getDate("2021-07-03"),
		},
		{
			ID:     "zogn",
			Amount: 86,
			Desc:   "E",
			Date:   getDate("2021-07-04"),
		},
		{
			ID:     "zojm",
			Amount: 77,
			Desc:   "F",
			Date:   getDate("2021-07-07"),
		},
		{
			ID:     "zoml",
			Amount: 65,
			Desc:   "G",
			Date:   getDate("2021-07-31"),
		},
		{
			ID:     "zopk",
			Amount: 66,
			Desc:   "H",
			Date:   getDate("2021-07-06"),
		},
		{
			ID:     "zosj",
			Amount: 56,
			Desc:   "I",
			Date:   getDate("2021-08-01"),
		},
		{
			ID:     "zovi",
			Amount: 73,
			Desc:   "J",
			Date:   getDate("2021-07-10"),
		},
	}

	if err := saveProxies(defaultProxies); err != nil {
		panic(err.Error())
	}
}

func PopulateSources() {
	defaultSources := []Source{
		{
			ID:     "zoUr",
			Amount: 24,
			Desc:   "A",
			Date:   getDate("2021-06-30"),
		},
		{
			ID:     "zoXq",
			Amount: 11,
			Desc:   "B",
			Date:   getDate("2021-06-30"),
		},
		{
			ID:     "zoap",
			Amount: 69,
			Desc:   "C",
			Date:   getDate("2021-07-01"),
		},
		{
			ID:     "zogn",
			Amount: 86,
			Desc:   "E",
			Date:   getDate("2021-07-04"),
		},
		{
			ID:     "zojm",
			Amount: 76,
			Desc:   "F",
			Date:   getDate("2021-07-07"),
		},
		{
			ID:     "zoml",
			Amount: 62,
			Desc:   "G",
			Date:   getDate("2021-07-31"),
		},
		{
			ID:     "zopk",
			Amount: 66,
			Desc:   "H",
			Date:   getDate("2021-07-06"),
		},
		{
			ID:     "zosj",
			Amount: 56,
			Desc:   "I",
			Date:   getDate("2021-08-01"),
		},
		{
			ID:     "zovi",
			Amount: 73,
			Desc:   "J",
			Date:   getDate("2021-07-10"),
		},
	}

	if err := saveSources(defaultSources); err != nil {
		panic(err.Error())
	}
}

func getDate(dateString string) time.Time {
	date, err := time.Parse("2006-01-02", dateString)

	if err != nil {
		panic(err.Error())
	}

	return date
}

func saveProxies(proxies []Proxy) error {
	tmpFile, err := ioutil.TempFile("", "proxy-sample-*.csv")

	if err != nil {
		return err
	}

	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)

	defer writer.Flush()

	CSVReport := NewCSVProxy(proxies)

	if err := writer.Write(CSVReport.Headers); err != nil {
		return err
	}

	if err := writer.WriteAll(CSVReport.Values); err != nil {
		return err
	}

	dst, err := os.Create(FileDataLocation + "proxy.csv")

	if err != nil {
		return err
	}

	defer dst.Close()

	input, err := ioutil.ReadFile(tmpFile.Name())

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst.Name(), input, 0644)

	if err != nil {
		return err
	}

	return nil
}

func saveSources(sources []Source) error {
	tmpFile, err := ioutil.TempFile("", "source-sample-*.csv")

	if err != nil {
		return err
	}

	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)

	defer writer.Flush()

	CSVReport := NewCSVSource(sources)

	if err := writer.Write(CSVReport.Headers); err != nil {
		return err
	}

	if err := writer.WriteAll(CSVReport.Values); err != nil {
		return err
	}

	dst, err := os.Create(FileDataLocation + "source.csv")

	if err != nil {
		return err
	}

	defer dst.Close()

	input, err := ioutil.ReadFile(tmpFile.Name())

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst.Name(), input, 0644)

	if err != nil {
		return err
	}

	return nil
}
