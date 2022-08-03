package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Reconcile() ([]Report, error) {
	proxies, err := db.FindAllProxies()

	if err != nil {
		return nil, err
	}

	var reports []Report

	for _, proxy := range proxies {
		source, err := db.FindSourceByID(proxy.ID)

		if err != nil {
			source = Source{}
		}

		reports = append(reports, Report{
			Proxy:   proxy,
			Source:  source,
			Remarks: getRemarks(proxy, source),
		})
	}

	return reports, nil
}

func ReconcileCSVFile(destinationdir string) error {
	reports, err := Reconcile()

	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile("", "reports-*.csv")

	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	writer := csv.NewWriter(tmpFile)

	defer writer.Flush()

	CSVReport := NewCSVReport(reports)

	if err := writer.Write(CSVReport.Headers); err != nil {
		return err
	}

	if err := writer.WriteAll(CSVReport.Values); err != nil {
		return err
	}

	dst, err := os.Create(destinationdir + "report.csv")

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

func ReconcileCSVFileByte() ([]byte, error) {
	var fileBytes []byte
	reports, err := Reconcile()

	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "reports-*.csv")

	if err != nil {
		return nil, err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	writer := csv.NewWriter(tmpFile)

	defer writer.Flush()

	CSVReport := NewCSVReport(reports)

	if err := writer.Write(CSVReport.Headers); err != nil {
		return nil, err
	}

	if err := writer.WriteAll(CSVReport.Values); err != nil {
		return nil, err
	}

	fileBytes, err = ioutil.ReadFile(tmpFile.Name())

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func ReconcileSummaryFile(destinationdir string) error {
	reports, err := Reconcile()

	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile("", "summary-*.txt")

	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	summaryReport := NewSummaryReport(reports)

	for _, v := range summaryReport.Headers {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " ")))

		if err != nil {
			return err
		}
	}

	tmpFile.WriteString("\n")
	tmpFile.WriteString("DISCREPANCIES \n")
	tmpFile.WriteString("============== \n")
	tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(summaryReport.DiscrepanciesHeaders, " | ")))

	for _, v := range summaryReport.DiscrepanciesValues {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " | ")))

		if err != nil {
			return err
		}
	}

	dst, err := os.Create(destinationdir + "summary.txt")

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

func ReconcileSummaryFileByte() ([]byte, error) {
	var fileBytes []byte
	reports, err := Reconcile()

	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "reports-*.txt")

	if err != nil {
		return nil, err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	summaryReport := NewSummaryReport(reports)

	for _, v := range summaryReport.Headers {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " ")))

		if err != nil {
			return nil, err
		}
	}

	tmpFile.WriteString("\n")
	tmpFile.WriteString("DISCREPANCIES \n")
	tmpFile.WriteString("============== \n")
	tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(summaryReport.DiscrepanciesHeaders, " | ")))

	for _, v := range summaryReport.DiscrepanciesValues {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " | ")))

		if err != nil {
			return nil, err
		}
	}

	fileBytes, err = ioutil.ReadFile(tmpFile.Name())

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func getRemarks(proxy Proxy, source Source) Remark {
	var remark Remark

	if source.ID == "" {
		remark.Discrepancies = append(remark.Discrepancies, "proxy data not found on source data")
		return remark
	}

	if proxy.Amount != source.Amount {
		remark.Discrepancies = append(remark.Discrepancies, fmt.Sprintf("proxy amount (%v) is different from source amount (%v)", proxy.Amount, source.Amount))
	}

	if proxy.Desc != source.Desc {
		remark.Discrepancies = append(remark.Discrepancies, fmt.Sprintf("proxy description (%s) is different from source description (%s)", proxy.Desc, source.Desc))
	}

	proxyDate := proxy.Date.Format("2006-01-02")
	sourceDate := source.Date.Format("2006-01-02")

	if proxyDate != sourceDate {
		remark.Discrepancies = append(remark.Discrepancies, fmt.Sprintf("proxy date (%s) is different from source date (%s)", proxyDate, sourceDate))
	}

	return remark
}
