package service

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database struct {
		Default string `yaml:"default"`
		Dsn     string `yaml:"dsn"`
		Name    string `yaml:"name"`
		User    string `yaml:"user"`
	} `yaml:"database"`
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
