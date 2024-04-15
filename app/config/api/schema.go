package api

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)
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

type Endpoint struct {
	Method string `yaml:"method"`
	Path   Path   `yaml:"path"`
}
type Path string
type Public struct {
	Path     Path `yaml:"path"`
	Services struct {
		Account  Account `yaml:"account"`
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
type Account struct {
	Endpoints struct {
		Create        Endpoint `yaml:"create"`
		CreateProfile Endpoint `yaml:"create_profile"`
		Login         Endpoint `yaml:"login"`
		Logout        Endpoint `yaml:"logout"`
		Whoami        Endpoint `yaml:"whoami"`
	} `yaml:"endpoints"`
	Path Path `yaml:"path"`
}
type Objectives struct {
	Endpoints struct {
		Attach    Endpoint `yaml:"attach"`
		Create    Endpoint `yaml:"create"`
		Delete    Endpoint `yaml:"delete"`
		Mark      Endpoint `yaml:"mark"`
		Placement Endpoint `yaml:"placement"`
	} `yaml:"endpoints"`
	Path Path `yaml:"path"`
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
