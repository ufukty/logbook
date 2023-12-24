package local

import "logbook/internal/web/serviced/models/services"

type IPAddress string

type Local map[services.ServiceName][]IPAddress

func (l *Local) ListPrivateIPs(service services.ServiceName) (ips []string) {
	if ipAddresses, ok := (*l)[service]; ok {
		for _, ipAddress := range ipAddresses {
			ips = append(ips, string(ipAddress))
		}
	}
	return
}
