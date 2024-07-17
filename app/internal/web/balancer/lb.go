package balancer

import (
	"fmt"
	"logbook/internal/utilities/randoms"
	"logbook/internal/web/logger"
	"logbook/models"

	"errors"
)

type InstanceSource interface {
	Instances() ([]models.Instance, error)
}

var (
	ErrNoHostAvailable = errors.New("no hosts are available right now")
)

type LoadBalancer struct {
	source InstanceSource
	index  int
}

var log = logger.NewLogger("LoadBalancer")

func New(is InstanceSource) *LoadBalancer {
	return &LoadBalancer{
		source: is,
		index:  0,
	}
}

func (lb *LoadBalancer) Next() (*models.Instance, error) {
	hosts, err := lb.source.Instances()
	if err != nil {
		return nil, fmt.Errorf("checking service pool: %w", err)
	}
	if len(hosts) == 0 {
		return nil, ErrNoHostAvailable
	}
	if len(hosts) <= lb.index {
		lb.index = randoms.UniformIntN(len(hosts))
	}
	var next = hosts[lb.index]
	lb.index = (lb.index + 1) % len(hosts)
	return &next, nil
}
