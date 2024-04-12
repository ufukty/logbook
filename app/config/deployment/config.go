package deployment

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Ports struct {
		Accounts   string `yaml:"accounts"`
		Objectives string `yaml:"objectives"`
	} `yaml:"ports"`
	Router struct {
		GracePeriod    time.Duration `yaml:"grace-period"`
		IdleTimeout    time.Duration `yaml:"idle-timeout"`
		ReadTimeout    time.Duration `yaml:"read-timeout"`
		RequestTimeout time.Duration `yaml:"request-timeout"`
		WriteTimeout   time.Duration `yaml:"write-timeout"`
	} `yaml:"router"`
	ServiceDiscovery struct {
		UpdatePeriod time.Duration `yaml:"update-period"`
	} `yaml:"service-discovery"`
}

func ReadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("opening config file: %w", err)
	}
	cfg := Config{}
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decoding config file: %w", err)
	}
	return cfg, nil
}
