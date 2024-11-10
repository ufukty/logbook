package balancer

import (
	"fmt"
	"logbook/models"
	"math/rand/v2"
	"sync"

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
	mu     sync.Mutex // index
}

func New(is InstanceSource) *LoadBalancer {
	return &LoadBalancer{
		source: is,
		index:  0,
	}
}

func (lb *LoadBalancer) Next() (*models.Instance, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	hosts, err := lb.source.Instances()
	if err != nil {
		return nil, fmt.Errorf("checking service pool: %w", err)
	}
	if len(hosts) == 0 {
		return nil, ErrNoHostAvailable
	}
	if len(hosts) <= lb.index {
		lb.index = rand.IntN(len(hosts))
	}
	var next = hosts[lb.index]
	lb.index = (lb.index + 1) % len(hosts)
	return &next, nil
}

func (lb *LoadBalancer) Host() (string, error) {
	h, err := lb.Next()
	if err != nil {
		return "", fmt.Errorf("next: %w", err)
	}
	return h.String(), nil
}

var _ Pool = &LoadBalancer{}
