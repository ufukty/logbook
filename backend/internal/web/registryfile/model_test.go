package registryfile

import (
	"fmt"
	"logbook/internal/startup"
	"testing"
)

func TestFileReader(t *testing.T) {
	l, deplcfg, _, err := startup.TestDependencies()
	if err != nil {
		t.Fatal(fmt.Errorf("startup: %w", err))
	}

	fr := NewFileReader("testdata/file.json", deplcfg, ServiceParams{8080, true}, l)
	instances, err := fr.Instances()
	defer fr.Stop()
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}
	if len(instances) != 2 {
		t.Error("assert")
	}
}
