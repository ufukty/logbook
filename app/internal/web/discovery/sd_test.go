package discovery

import (
	"fmt"
	"logbook/models"
	"testing"
	"time"
)

func TestServiceDiscoveryStage(t *testing.T) {
	var tcs = map[models.Environment]string{
		"local": "models/local/service_discovery.yml",
		"stage": "models/stage/service_discovery.json",
	}
	for env, tc := range tcs {
		t.Run(string(env), func(t *testing.T) {
			sd := New(env, tc, time.Second*5)
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
