package models

import (
	"logbook/internal/web/serviced/models/provider"
	"logbook/internal/web/serviced/models/provider/digitalocean"
	"logbook/internal/web/serviced/models/provider/local"
	"logbook/internal/web/serviced/models/services"
)

type ServiceDiscoveryFile struct {
	Digitalocean digitalocean.Digitalocean `json:"digitalocean"`
	Local        local.Local               `json:"local"`
}

func (f ServiceDiscoveryFile) ListPrivateIPs(service services.ServiceName) (ips []string) {
	var providers = []provider.Provider{&f.Digitalocean, &f.Local}
	for _, provider := range providers {
		ips = append(ips, provider.ListPrivateIPs(service)...)
	}
	return
}
