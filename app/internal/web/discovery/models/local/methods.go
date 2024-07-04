package local

import (
	"fmt"
	"logbook/models"
	"slices"
)

func (service Accounts) ListPrivateIPs() []string {
	var ips = []string{}
	for _, vps := range service {
		ips = append(ips, vps)
	}
	return ips
}

func (service Objectives) ListPrivateIPs() []string {
	var ips = []string{}
	for _, vps := range service {
		ips = append(ips, vps)
	}
	return ips
}

func (service ApiGateway) ListPrivateIPs() []string {
	var ips = []string{}
	for _, vps := range service {
		ips = append(ips, vps)
	}
	return ips
}

func (l Local) ServicePool(service models.Service) ([]string, error) {
	switch service {
	case models.Account:
		return slices.Concat(
			l.Accounts.ListPrivateIPs(),
		), nil
	case models.Gateway:
		return slices.Concat(
			l.ApiGateway.ListPrivateIPs(),
		), nil
	case models.Objectives:
		return slices.Concat(
			l.Objectives.ListPrivateIPs(),
		), nil
	default:
		return nil, fmt.Errorf("unrecognized service name %q", service)
	}
}
