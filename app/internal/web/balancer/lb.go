package balancer

import (
	"fmt"
	"logbook/internal/utilities/randoms"
	"logbook/internal/web/discovery"
	"logbook/internal/web/logger"
	"logbook/models"

	"errors"
)

var (
	ErrNoHostAvailable = errors.New("no hosts are available right now")
)

type LoadBalancer struct {
	sd            *discovery.ServiceDiscovery
	index         int
	targetService models.Service
}

var log = logger.NewLogger("LoadBalancer")

// service: ip address and host
// hosts: ip addresses of available hosts
// port: port which will be used as forwarded target
func New(sd *discovery.ServiceDiscovery, targetService models.Service) *LoadBalancer {
	return &LoadBalancer{
		sd:            sd,
		index:         0,
		targetService: targetService,
	}
}

func (lb *LoadBalancer) Next() (string, error) {
	hosts, err := lb.sd.ServicePool(lb.targetService)
	if err != nil {
		return "", fmt.Errorf("checking service pool: %w", err)
	}
	if len(hosts) == 0 {
		return "", ErrNoHostAvailable
	}
	if len(hosts) <= lb.index {
		lb.index = randoms.UniformIntN(len(hosts))
	}
	var next = hosts[lb.index]
	lb.index = (lb.index + 1) % len(hosts)
	return next, nil
}
