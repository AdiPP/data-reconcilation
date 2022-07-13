package mapper

import (
	"testing"

	"github.com/AdiPP/recon-general/loader"
)

func TestMapProxyFile(t *testing.T) {
	file, _ := loader.NewLoader().Load(proxy_path)
	_, err := NewSourceMapper().Map(file)

	if err != nil {
		t.Errorf(err.Error())
	}

	file.Close()
}
