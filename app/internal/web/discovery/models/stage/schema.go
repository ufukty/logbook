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
	Fra1 Fra1 `json:"fra1"`
}
type Droplet struct {
	Backups            bool     `json:"backups"`
	CreatedAt          string   `json:"created_at"`
	Disk               float64  `json:"disk"`
	DropletAgent       any      `json:"droplet_agent"`
	GracefulShutdown   bool     `json:"graceful_shutdown"`
	Id                 string   `json:"id"`
	Image              string   `json:"image"`
	Ipv4Address        string   `json:"ipv4_address"`
	Ipv4AddressPrivate string   `json:"ipv4_address_private"`
	Ipv6               bool     `json:"ipv6"`
	Ipv6Address        string   `json:"ipv6_address"`
	Locked             bool     `json:"locked"`
	Memory             float64  `json:"memory"`
	Monitoring         bool     `json:"monitoring"`
	Name               string   `json:"name"`
	PriceHourly        float64  `json:"price_hourly"`
	PriceMonthly       float64  `json:"price_monthly"`
	PrivateNetworking  bool     `json:"private_networking"`
	Region             string   `json:"region"`
	ResizeDisk         bool     `json:"resize_disk"`
	Size               string   `json:"size"`
	SshKeys            []string `json:"ssh_keys"`
	Status             string   `json:"status"`
	Tags               []string `json:"tags"`
	Timeouts           any      `json:"timeouts"`
	Urn                string   `json:"urn"`
	UserData           any      `json:"user_data"`
	Vcpus              float64  `json:"vcpus"`
	VolumeIds          []any    `json:"volume_ids"`
	VpcUuid            string   `json:"vpc_uuid"`
}
type Fra1 struct {
	Services struct {
		Account    Account    `json:"account"`
		Gateway    Gateway    `json:"gateway"`
		Objectives Objectives `json:"objectives"`
	} `json:"services"`
	Vpc struct {
		CreatedAt   string `json:"created_at"`
		Default     bool   `json:"default"`
		Description string `json:"description"`
		Id          string `json:"id"`
		IpRange     string `json:"ip_range"`
		Name        string `json:"name"`
		Region      string `json:"region"`
		Urn         string `json:"urn"`
	} `json:"vpc"`
}
type Gateway []Droplet
type Objectives []Droplet
type Config struct {
	Digitalocean DigitalOcean `json:"digitalocean"`
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
