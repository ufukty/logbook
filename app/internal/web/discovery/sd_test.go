package discovery

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func TestServiceDiscoveryStage(t *testing.T) {
	var tcs = []string{
		"service_discovery_local.json",
		"service_discovery_stage.json",
	}
	for _, tc := range tcs {
		t.Run(tc, func(t *testing.T) {
			sd := New(filepath.Join("testdata", tc), time.Second*5)
			ips, err := sd.ServicePool("objectives")
			if err != nil {
				t.Fatal(fmt.Errorf("act: %w", err))
			}
			if len(ips) != 1 {
				t.Error("assert, 1")
			}
		})
	}

}
