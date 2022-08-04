package exporting

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/AdiPP/reconciliation/pkg/listing"
)

const (
	TypeProxyNotFound   = "PROXY_NOT_FOUND"
	TypeDiffProxyAmount = "DIFFERENT_PROXY_AMOUNT"
	TypeDiffProxyDesc   = "DIFFERENT_PROXY_DESCRIPTION"
	TypeDiffProxyDate   = "DIFFERENT_PROXY_DATE"
)

var (
	ProxyNotFound   = "proxy data not found on source data"
	DiffProxyAmount = "proxy amount (%v) is different from source amount (%v)"
	DiffProxyDesc   = "proxy description (%s) is different from source description (%s)"
	DiffProxyDate   = "proxy date (%s) is different from source date (%s)"
)

type ReportOption struct {
	StartDate time.Time
	EndDate   time.Time
}

type Repository interface {
	FindAllProxies() ([]listing.Proxy, error)
	FindSourceByID(ID string) (listing.Source, error)
}

type Service interface {
	FindSourceByID(ID string) (Source, error)
	GetReportData(reportOpt ReportOption) ([]Report, error)
	GenerateReportFile(reportOpt ReportOption, dstDir string) error
	GenerateSummaryReportFile(reportOpt ReportOption, dstDir string) error
	WriteReportFile(writer *csv.Writer, reports []Report) error
	WriteSummaryReportFile(reportOpt ReportOption, tmpFile *os.File, reports []Report) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) FindSourceByID(ID string) (Source, error) {
	var source Source

	storageSource, err := s.r.FindSourceByID(ID)

	if err != nil {
		return source, err
	}

	return Source{
		ID:     storageSource.ID,
		Amount: storageSource.Amount,
		Desc:   storageSource.Desc,
		Date:   storageSource.Date,
	}, nil
}

func (s *service) GetReportData(reportOpt ReportOption) ([]Report, error) {
	storageProxies, err := s.r.FindAllProxies()

	if err != nil {
		return nil, err
	}

	var proxies []Proxy

	for _, v := range storageProxies {
		proxies = append(
			proxies, Proxy{
				ID:     v.ID,
				Amount: v.Amount,
				Desc:   v.Desc,
				Date:   v.Date,
			},
		)
	}

	filteredProxies := getFilteredProxies(reportOpt, proxies)
	sortedProxies := getSorteredProxies(filteredProxies)

	var reports []Report

	for _, v := range sortedProxies {
		source, err := s.FindSourceByID(v.ID)

		if err != nil {
			if err != listing.ErrSourceNotFound {
				return nil, err
			}
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

	err = s.WriteSummaryReportFile(reportOpt, tmpFile, reports)

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

func (s *service) WriteSummaryReportFile(reportOpt ReportOption, tmpFile *os.File, reports []Report) error {
	export := NewTextSummaryReportExport(reportOpt, reports)

	tmpFile.WriteString("# RECONCILIATION REPORT SUMMARY \n")
	tmpFile.WriteString("## QUERY \n")

	for _, v := range export.Queries {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " ")))

		if err != nil {
			return err
		}
	}

	tmpFile.WriteString("\n")
	tmpFile.WriteString("## SUMMARY \n")

	for _, v := range export.Summaries {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " ")))

		if err != nil {
			return err
		}
	}

	tmpFile.WriteString("\n")
	tmpFile.WriteString("## DISCREPANCIES \n")

	for _, v := range export.Discrepancies {
		_, err := tmpFile.WriteString(fmt.Sprintf("%s \n", strings.Join(v, " ")))

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

func getSorteredProxies(proxies []Proxy) []Proxy {
	sort.Slice(proxies, func(i, j int) bool {
		return proxies[i].Date.Before(proxies[j].Date)
	})

	return proxies
}

func getReportRemarks(proxy Proxy, source Source) []Remark {
	var remarks []Remark

	if source.ID == "" {
		remarks = append(
			remarks,
			Remark{
				Type:  TypeProxyNotFound,
				Error: fmt.Errorf(ProxyNotFound),
			},
		)

		return remarks
	}

	if proxy.Amount != source.Amount {
		remarks = append(
			remarks,
			Remark{
				Type:  TypeDiffProxyAmount,
				Error: fmt.Errorf(DiffProxyAmount, proxy.Amount, source.Amount),
			},
		)
	}

	if proxy.Desc != source.Desc {
		remarks = append(
			remarks,
			Remark{
				Type:  TypeDiffProxyDesc,
				Error: fmt.Errorf(DiffProxyDesc, proxy.Desc, source.Desc),
			},
		)
	}

	proxyDate := proxy.Date.Format("2006-01-02")
	sourceDate := source.Date.Format("2006-01-02")

	if proxyDate != sourceDate {
		remarks = append(
			remarks,
			Remark{
				Type:  TypeDiffProxyDate,
				Error: fmt.Errorf(DiffProxyDate, proxyDate, sourceDate),
			},
		)
	}

	return remarks
}
