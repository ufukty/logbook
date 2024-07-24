package registryfile

import (
	"fmt"
	"testing"
	"time"
)

func TestFileReader(t *testing.T) {
	fr := NewFileReader("testdata/file.json", time.Second, ServiceParams{8080, true})
	instances, err := fr.Instances()
	defer fr.Stop()
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}
	if len(instances) != 2 {
		t.Error("assert")
	}
}
