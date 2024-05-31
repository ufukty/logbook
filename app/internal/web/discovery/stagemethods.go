package discovery

import (
	"fmt"
	"logbook/models"
	"slices"
)

func (service Account) ListPrivateIPs() []string {
	var ips = []string{}
	for _, vps := range service {
		ips = append(ips, vps.Ipv4AddressPrivate)
	}
	return ips
}

func (service Objectives) ListPrivateIPs() []string {
	var ips = []string{}
	for _, vps := range service {
		ips = append(ips, vps.Ipv4AddressPrivate)
	}
	return ips
}

func (service Gateway) ListPrivateIPs() []string {
	var ips = []string{}
	for _, vps := range service {
		ips = append(ips, vps.Ipv4AddressPrivate)
	}
	return ips
}

func (s Stage) ServicePool(service models.Service) ([]string, error) {
	switch service {
	case models.Account:
		return slices.Concat(
			s.Digitalocean.Fra1.Services.Account.ListPrivateIPs(),
		), nil
	case models.Gateway:
		return slices.Concat(
			s.Digitalocean.Fra1.Services.Gateway.ListPrivateIPs(),
		), nil
	case models.Objectives:
		return slices.Concat(
			s.Digitalocean.Fra1.Services.Objectives.ListPrivateIPs(),
		), nil
	default:
		return nil, fmt.Errorf("unrecognized service name %q", service)
	}
}
