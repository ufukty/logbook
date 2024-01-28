package config

import (
	"logbook/internal/utilities/reflux"
	"testing"
)

func Test_ReadConfig(t *testing.T) {
	cfg := Read("testdata/config.yml")
	reflux.Print(cfg)

	if cfg.Tasks.ServiceDiscoveryConfig != "75c31fcc-6dca-5e99-9bad-ea82ad9fe1f6" {
		t.Error("validation")
	}
}
