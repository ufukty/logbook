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
	case "gateway":
		return slices.Concat(
			s.Digitalocean.Fra1.Services.Gateway.ListPrivateIPs(),
		), nil
	case "objectives":
		return slices.Concat(
			s.Digitalocean.Fra1.Services.Objectives.ListPrivateIPs(),
		), nil
	case "account":
		return slices.Concat(
			s.Digitalocean.Fra1.Services.Account.ListPrivateIPs(),
		), nil
	default:
		return nil, fmt.Errorf("unrecognized service name %q", service)
	}
}
