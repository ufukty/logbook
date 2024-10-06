package registryfile

import (
	"fmt"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"testing"
	"time"
)

func TestFileReader(t *testing.T) {
	deplycfg := &deployment.Config{
		RegistryFile: struct {
			UpdatePeriod time.Duration "yaml:\"update-period\""
		}{
			UpdatePeriod: time.Second,
		},
	}
	fr := NewFileReader("testdata/file.json", deplycfg, ServiceParams{8080, true}, logger.New("test"))
	instances, err := fr.Instances()
	defer fr.Stop()
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}
	if len(instances) != 2 {
		t.Error("assert")
	}
}
