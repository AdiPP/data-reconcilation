package mapper

import (
	"testing"

	"github.com/AdiPP/recon-general/loader"
)

func TestMapProxyFile(t *testing.T) {
	proxyFilePath := "../source_test/valid_proxy_test.csv"
	file, _ := loader.NewLoader().Load(proxyFilePath)
	_, err := NewProxyMapper().Map(file)

	if err != nil {
		t.Errorf(err.Error())
	}

	file.Close()
}

func TestInvalidHeaderProxyFile(t *testing.T) {
	proxyFilePath := "../source_test/invalid_header_proxy_test.csv"
	file, _ := loader.NewLoader().Load(proxyFilePath)
	_, err := NewProxyMapper().Map(file)
	expectation := FILE_NOT_SUPPORTED

	if err.Error() != expectation {
		t.Errorf("Expected erorr %v but got %v", expectation, err.Error())
	}

	file.Close()
}
