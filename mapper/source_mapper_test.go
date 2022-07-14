package mapper

import (
	"testing"

	"github.com/AdiPP/recon-general/loader"
)

func TestMapSourceFile(t *testing.T) {
	sourceFilePath := "../source_test/valid_source_test.csv"
	file, _ := loader.NewLoader().Load(sourceFilePath)
	_, err := NewSourceMapper().Map(file)

	if err != nil {
		t.Errorf(err.Error())
	}

	file.Close()
}

func TestInvalidHeaderSourceFile(t *testing.T) {
	proxyFilePath := "../source_test/invalid_header_source_test.csv"
	file, _ := loader.NewLoader().Load(proxyFilePath)
	_, err := NewSourceMapper().Map(file)
	expectation := FILE_NOT_SUPPORTED

	if err.Error() != expectation {
		t.Errorf("Expected erorr %v but got %v", expectation, err.Error())
	}

	file.Close()
}
