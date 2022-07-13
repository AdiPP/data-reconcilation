package processor

import (
	"testing"
)

const (
	source_path = "../source_test/source_test.csv"
	proxy_path  = "../source_test/proxy_test.csv"
)

func TestProcessReport(t *testing.T) {
	actual := len(Process(source_path, proxy_path))
	expectation := 10

	if actual != expectation {
		t.Errorf("Expected %v but got %v", expectation, actual)
	}
}

func TestProcessCSV(t *testing.T) {
	ProcessCSV(source_path, proxy_path)
}

func TestProcessSummaryReport(t *testing.T) {
	ProcessSummary(source_path, proxy_path)
}
