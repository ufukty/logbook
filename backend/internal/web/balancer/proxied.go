package balancer

import (
	"fmt"
	"logbook/internal/utils/urls"
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
	return urls.Join(h, p.Suffix), nil
}

var _ Pool = Proxied{}
