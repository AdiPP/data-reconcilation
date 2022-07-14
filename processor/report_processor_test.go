package processor

import (
	"testing"

	"github.com/AdiPP/recon-general/mapper"
)

func TestSuccessProcessReport(t *testing.T) {
	source_path := "../source_test/valid_source_test.csv"
	proxy_path := "../source_test/valid_proxy_test.csv"
	results, err := Process(source_path, proxy_path)

	if err != nil {
		t.Errorf(err.Error())
	}

	actual := len(results)
	expectation := 10

	if actual != expectation {
		t.Errorf("Expected %v but got %v", expectation, actual)
	}
}

func TestFailedProcessReportCausedByInvalidFile(t *testing.T) {
	source_path := "../source_test/invalid_source_test.csv"
	proxy_path := "../source_test/invalid_proxy_test.csv"
	_, err := Process(source_path, proxy_path)

	expectation := mapper.FILE_NOT_SUPPORTED

	if err == nil {
		t.Errorf("Expected erorr %v but got %v", expectation, err.Error())
	}
}
