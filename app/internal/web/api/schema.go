package api

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Method string `yaml:"method"`
	Path   Path   `yaml:"path"`
}
type Path string
type Public struct {
	Path     Path `yaml:"path"`
	Services struct {
		Document struct {
			Endpoints struct {
				List Endpoint `yaml:"list"`
			} `yaml:"endpoints"`
			Path Path `yaml:"path"`
		} `yaml:"document"`
		Objectives Objectives `yaml:"objectives"`
		Tags       struct {
			Endpoints struct {
				Assign   Endpoint `yaml:"assign"`
				Creation Endpoint `yaml:"creation"`
			} `yaml:"endpoints"`
			Path Path `yaml:"path"`
		} `yaml:"tags"`
	} `yaml:"services"`
}
type Objectives struct {
	Endpoints struct {
		Create       Endpoint `yaml:"create"`
		GetPlacement Endpoint `yaml:"getPlacement"`
	} `yaml:"endpoints"`
	Path Path `yaml:"path"`
}
// IMPORTANT:
// Types are defined only for internal purposes.
// Do not refer auto generated type names from outside.
// Because they will change as config schema changes.
type autoGenA struct {
	Public Public `yaml:"public"`
}

func (a autoGenA) Range() map[string]Public {
	return map[string]Public{"public": a.Public}
}

type Config struct {
	Domain   string   `yaml:"domain"`
	Gateways autoGenA `yaml:"gateways"`
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
