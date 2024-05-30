package stage

import (
	"encoding/json"
	"fmt"
	"os"
)

func (a Config) Range() map[string]DigitalOcean {
	return map[string]DigitalOcean{"digitalocean": a.Digitalocean}
}

type Account []Droplet
type DigitalOcean struct {
	Fra1 Fra1 `yaml:"fra1"`
}
type Droplet struct {
	Backups            bool     `yaml:"backups"`
	CreatedAt          string   `yaml:"created_at"`
	Disk               float64  `yaml:"disk"`
	DropletAgent       any      `yaml:"droplet_agent"`
	GracefulShutdown   bool     `yaml:"graceful_shutdown"`
	Id                 string   `yaml:"id"`
	Image              string   `yaml:"image"`
	Ipv4Address        string   `yaml:"ipv4_address"`
	Ipv4AddressPrivate string   `yaml:"ipv4_address_private"`
	Ipv6               bool     `yaml:"ipv6"`
	Ipv6Address        string   `yaml:"ipv6_address"`
	Locked             bool     `yaml:"locked"`
	Memory             float64  `yaml:"memory"`
	Monitoring         bool     `yaml:"monitoring"`
	Name               string   `yaml:"name"`
	PriceHourly        float64  `yaml:"price_hourly"`
	PriceMonthly       float64  `yaml:"price_monthly"`
	PrivateNetworking  bool     `yaml:"private_networking"`
	Region             string   `yaml:"region"`
	ResizeDisk         bool     `yaml:"resize_disk"`
	Size               string   `yaml:"size"`
	SshKeys            []string `yaml:"ssh_keys"`
	Status             string   `yaml:"status"`
	Tags               []string `yaml:"tags"`
	Timeouts           any      `yaml:"timeouts"`
	Urn                string   `yaml:"urn"`
	UserData           any      `yaml:"user_data"`
	Vcpus              float64  `yaml:"vcpus"`
	VolumeIds          []any    `yaml:"volume_ids"`
	VpcUuid            string   `yaml:"vpc_uuid"`
}
type Fra1 struct {
	Services struct {
		Account    Account    `yaml:"account"`
		Gateway    Gateway    `yaml:"gateway"`
		Objectives Objectives `yaml:"objectives"`
	} `yaml:"services"`
	Vpc struct {
		CreatedAt   string `yaml:"created_at"`
		Default     bool   `yaml:"default"`
		Description string `yaml:"description"`
		Id          string `yaml:"id"`
		IpRange     string `yaml:"ip_range"`
		Name        string `yaml:"name"`
		Region      string `yaml:"region"`
		Urn         string `yaml:"urn"`
	} `yaml:"vpc"`
}
type Gateway []Droplet
type Objectives []Droplet
type Config struct {
	Digitalocean DigitalOcean `yaml:"digitalocean"`
}

func ReadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("opening config file: %w", err)
	}
	cfg := Config{}
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decoding config file: %w", err)
	}
	return cfg, nil
}
