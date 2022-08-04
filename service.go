package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type ReconcileOption struct {
	StartDate time.Time
	EndDate   time.Time
}

func Reconcile(opt ReconcileOption) ([]Report, error) {
	proxies, err := getFilteredProxies(opt)

	if err != nil {
		return nil, err
	}

	var reports []Report

	for _, v := range proxies {
		source, err := db.FindSourceByID(v.ID)

		if err != nil {
			source = Source{}
		}

		reports = append(reports, Report{
			Proxy:   v,
			Source:  source,
			Remarks: getRemarks(v, source),
		})
	}

	return reports, nil
}

func ReconcileReportFile(opt ReconcileOption, destinationdir string) error {
	var err error

	reports, err := Reconcile(opt)

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

	err = writeReconcileReportFile(writer, reports)

	if err != nil {
		return err
	}

	now := time.Now()
	filename := fmt.Sprintf(
		"reconciliation-report-%s%s.csv",
		now.Format("20060102"),
		now.Format("150405"),
	)
	dstFile, err := os.Create(destinationdir + filename)

	if err != nil {
		return err
	}

	defer dstFile.Close()

	err = copyTempFileToDestinationFile(tmpFile, dstFile)

	if err != nil {
		return err
	}

	return nil
}

func ReconcileReportFileByte(opt ReconcileOption) ([]byte, error) {
	var fileBytes []byte
	var err error

	reports, err := Reconcile(opt)

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

	err = writeReconcileReportFile(writer, reports)

	if err != nil {
		return nil, err
	}

	fileBytes, err = ioutil.ReadFile(tmpFile.Name())

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func ReconcileReportSummaryFile(opt ReconcileOption, destinationdir string) error {
	reports, err := Reconcile(opt)

	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile("", "summary-*.txt")

	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = writeReconcileSummaryFile(tmpFile, reports)

	if err != nil {
		return err
	}

	now := time.Now()
	filename := fmt.Sprintf(
		"reconciliation-report-summary-%s%s.txt",
		now.Format("20060102"),
		now.Format("150405"),
	)
	dstFile, err := os.Create(destinationdir + filename)

	if err != nil {
		return err
	}

	defer dstFile.Close()

	err = copyTempFileToDestinationFile(tmpFile, dstFile)

	if err != nil {
		return err
	}

	return nil
}

func ReconcileReportSummaryFileByte(opt ReconcileOption) ([]byte, error) {
	var fileBytes []byte
	reports, err := Reconcile(opt)

	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "reports-*.txt")

	if err != nil {
		return nil, err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = writeReconcileSummaryFile(tmpFile, reports)

	if err != nil {
		return nil, err
	}

	fileBytes, err = ioutil.ReadFile(tmpFile.Name())

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func getFilteredProxies(opt ReconcileOption) ([]Proxy, error) {
	proxies, err := db.FindAllProxies()

	if err != nil {
		return nil, err
	}

	if !opt.StartDate.IsZero() || !opt.EndDate.IsZero() {
		var filteredProxies []Proxy

		for _, v := range proxies {
			if v.Date == opt.StartDate || v.Date == opt.EndDate || (v.Date.After(opt.StartDate) && v.Date.Before(opt.EndDate)) {
				filteredProxies = append(filteredProxies, v)
			}
		}

		return filteredProxies, nil
	}

	return proxies, nil
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

func writeReconcileReportFile(writer *csv.Writer, reports []Report) error {
	CSVReport := NewCSVReport(reports)

	if err := writer.Write(CSVReport.Headers); err != nil {
		return err
	}

	if err := writer.WriteAll(CSVReport.Values); err != nil {
		return err
	}

	return nil
}

func writeReconcileSummaryFile(tmpFile *os.File, reports []Report) error {
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

	return nil
}

func copyTempFileToDestinationFile(tmpFile *os.File, dstFile *os.File) error {
	input, err := ioutil.ReadFile(tmpFile.Name())

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dstFile.Name(), input, 0644)

	if err != nil {
		return err
	}

	return nil
}
