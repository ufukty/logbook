package config

import (
	"logbook/config/reader"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func Test_ReadConfig(t *testing.T) {
	os.Args = []string{os.Args[0], "-config", "testdata/config.yml"}

	var config = Config{}
	reader.PopulateConfig(&config)
	if config.Tasks.ServiceDiscoveryConfig != "75c31fcc-6dca-5e99-9bad-ea82ad9fe1f6" {
		t.Error("validation")
	}

	yaml.NewEncoder(os.Stdout).Encode(config)
}
