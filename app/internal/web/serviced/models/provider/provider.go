package provider

import "logbook/internal/web/serviced/models/services"

type Provider interface {
	ListPrivateIPs(services.ServiceName) []string
}
