package mapper

import (
	"testing"

	"github.com/AdiPP/recon-general/loader"
)

func TestMapSourceFile(t *testing.T) {
	file, _ := loader.NewLoader().Load(source_path)
	_, err := NewProxyMapper().Map(file)

	if err != nil {
		t.Errorf(err.Error())
	}

	file.Close()
}
