package service

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
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
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()
	c := Config{}
	err = yaml.NewDecoder(file).Decode(&c)
	if err != nil {
		return Config{}, fmt.Errorf("decoding config file: %w", err)
	}
	return c, nil
}
