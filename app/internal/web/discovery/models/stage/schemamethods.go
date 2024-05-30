package stage

import "fmt"

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

func (config Config) ServicePool(service string) ([]string, error) {
	switch service {
	case "gateway":
		return config.Digitalocean.Fra1.Services.Gateway.ListPrivateIPs(), nil
	case "objectives":
		return config.Digitalocean.Fra1.Services.Objectives.ListPrivateIPs(), nil
	case "account":
		return config.Digitalocean.Fra1.Services.Account.ListPrivateIPs(), nil
	default:
		return nil, fmt.Errorf("unrecognized service name %q", service)
	}
}
