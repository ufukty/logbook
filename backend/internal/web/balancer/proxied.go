package balancer

import (
	"fmt"
	"logbook/models"
)

type Proxied struct {
	Lb     *LoadBalancer
	Suffix string
}

func NewProxied(is InstanceSource, suffix models.Service) Proxied {
	return Proxied{
		Lb:     New(is),
		Suffix: string(suffix),
	}
}

func (p Proxied) Host() (string, error) {
	h, err := p.Lb.Host()
	if err != nil {
		return "", fmt.Errorf("Lb.Host: %w", err)
	}
	return h + p.Suffix, nil
}

var _ Pool = Proxied{}
