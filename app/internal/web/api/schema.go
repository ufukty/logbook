package api

import (
	"fmt"
	"logbook/internal/web/paths"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Domain   string `yaml:"Domain"`
	Gateways struct {
		Public struct {
			Path     string `yaml:"Path"`
			Services struct {
				Document struct {
					Path      string `yaml:"Path"`
					Endpoints struct {
						List paths.Endpoint `yaml:"List"`
					} `yaml:"Endpoints"`
				} `yaml:"Document"`
				Objectives struct {
					Path      string `yaml:"Path"`
					Endpoints struct {
						Create struct {
							Path   string `yaml:"Path"`
							Method string `yaml:"Method"`
						} `yaml:"Create"`
						GetPlacement paths.Endpoint `yaml:"GetPlacement"`
					} `yaml:"Endpoints"`
				} `yaml:"Objectives"`
				Tags struct {
					Path      string `yaml:"Path"`
					Endpoints struct {
						Creation paths.Endpoint `yaml:"Creation"`
						Assign   paths.Endpoint `yaml:"Assign"`
					} `yaml:"Endpoints"`
				} `yaml:"Tags"`
			} `yaml:"Services"`
		} `yaml:"Public"`
	} `yaml:"Gateways"`
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
