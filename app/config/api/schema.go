package api

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Account struct {
	Endpoints struct {
		Create        Endpoint `yaml:"create"`
		CreateProfile Endpoint `yaml:"create_profile"`
		Login         Endpoint `yaml:"login"`
		Logout        Endpoint `yaml:"logout"`
		Whoami        Endpoint `yaml:"whoami"`
	} `yaml:"endpoints"`
	Path Path `yaml:"path"`
	Port Port `yaml:"port"`
}
type Document struct {
	Endpoints struct {
		List Endpoint `yaml:"list"`
	} `yaml:"endpoints"`
	Path Path `yaml:"path"`
	Port Port `yaml:"port"`
}
type Endpoint struct {
	Method string `yaml:"method"`
	Path   Path   `yaml:"path"`
}
type Gateway struct {
	Path Path `yaml:"path"`
	Port Port `yaml:"port"`
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
	Port Port `yaml:"port"`
}
type Path string
type Port int
type Tags struct {
	Endpoints struct {
		Assign   Endpoint `yaml:"assign"`
		Creation Endpoint `yaml:"creation"`
	} `yaml:"endpoints"`
	Path Path `yaml:"path"`
	Port Port `yaml:"port"`
}
type Config struct {
	Account    Account    `yaml:"account"`
	Document   Document   `yaml:"document"`
	Gateway    Gateway    `yaml:"gateway"`
	Objectives Objectives `yaml:"objectives"`
	Tags       Tags       `yaml:"tags"`
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
