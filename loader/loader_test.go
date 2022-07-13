package loader

import "testing"

const (
	source_path = "../source/proxy.csv"
)

func TestLoadFile(t *testing.T) {
	loader := Loader{}
	file, err := loader.Load(source_path)

	if err != nil {
		t.Errorf(err.Error())
	}

	file.Close()
}
