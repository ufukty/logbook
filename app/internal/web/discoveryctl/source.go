package discoveryctl

import (
	"logbook/internal/web/balancer"
	"logbook/models"
)

type source struct {
	s models.Service
	c *Client
}

var _ balancer.InstanceSource = &source{}

func newServiceStore(c *Client, service models.Service) *source {
	return &source{
		s: service,
		c: c,
	}
}

func (d *source) Instances() ([]models.Instance, error) {
	return d.c.instances(d.s)
}
