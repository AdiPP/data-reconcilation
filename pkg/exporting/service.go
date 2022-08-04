package exporting

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Todo-Adi: Add errors const

type ReportOption struct {
	StartDate time.Time
	EndDate   time.Time
}

type Repository interface {
	FindAllProxies() ([]Proxy, error)
	FindSourceByID(ID string) (Source, error)
}

type Service interface {
	GetReportData(reportOpt ReportOption) ([]Report, error)
	GenerateReportFile(reportOpt ReportOption, dstDir string) error
	GenerateSummaryReportFile(reportOpt ReportOption, dstDir string) error
	WriteReportFile(writer *csv.Writer, reports []Report) error
	WriteSummaryReportFile(tmpFile *os.File, reports []Report) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetReportData(reportOpt ReportOption) ([]Report, error) {
	proxies, err := s.r.FindAllProxies()
	if err != nil {
		return nil, err
	}

	filteredProxies := getFilteredProxies(reportOpt, proxies)

	var reports []Report

	for _, v := range filteredProxies {
		source, err := s.r.FindSourceByID(v.ID)

		if err != nil {
			source = Source{}
		}

		reports = append(reports, Report{
			Proxy:   v,
			Source:  source,
			Remarks: getReportRemarks(v, source),
		})
	}

	return reports, nil
}

func (s *service) GenerateReportFile(reportOpt ReportOption, dstDir string) error {
	reports, err := s.GetReportData(reportOpt)

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

	err = s.WriteReportFile(writer, reports)

	if err != nil {
		return err
	}

	now := time.Now()
	filename := fmt.Sprintf(
		"reconciliation-report-%s%s.csv",
		now.Format("20060102"),
		now.Format("150405"),
	)

	dstFile, err := os.Create(dstDir + filename)

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

func (s *service) GenerateSummaryReportFile(reportOpt ReportOption, dstDir string) error {
	reports, err := s.GetReportData(reportOpt)

	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile("", "summary-*.txt")

	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = s.WriteSummaryReportFile(tmpFile, reports)

	if err != nil {
		return err
	}

	now := time.Now()
	filename := fmt.Sprintf(
		"reconciliation-report-summary-%s%s.txt",
		now.Format("20060102"),
		now.Format("150405"),
	)
	dstFile, err := os.Create(dstDir + filename)

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

func getFilteredProxies(reportOpt ReportOption, proxies []Proxy) []Proxy {

	if !reportOpt.StartDate.IsZero() || !reportOpt.EndDate.IsZero() {
		var filteredProxies []Proxy

		for _, v := range proxies {
			if v.Date == reportOpt.StartDate || v.Date == reportOpt.EndDate || (v.Date.After(reportOpt.StartDate) && v.Date.Before(reportOpt.EndDate)) {
				filteredProxies = append(filteredProxies, v)
			}
		}

		return filteredProxies
	}

	return proxies
}

func getReportRemarks(proxy Proxy, source Source) Remark {
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

func (s *service) WriteReportFile(writer *csv.Writer, reports []Report) error {
	export := NewCSVReportExport(reports)

	if err := writer.Write(export.Headers); err != nil {
		return err
	}

	if err := writer.WriteAll(export.Values); err != nil {
		return err
	}

	return nil
}

func (s *service) WriteSummaryReportFile(tmpFile *os.File, reports []Report) error {
	export := NewTextSummaryReportExport(reports)

	for _, v := range export.Headers {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " ")))

		if err != nil {
			return err
		}
	}

	tmpFile.WriteString("\n")
	tmpFile.WriteString("DISCREPANCIES \n")
	tmpFile.WriteString("============== \n")
	tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(export.DiscrepanciesHeaders, " | ")))

	for _, v := range export.DiscrepanciesValues {
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
