package sidecar

import (
	"logbook/internal/web/balancer"
	"logbook/models"
)

type source struct {
	s models.Service
	c *Sidecar
}

var _ balancer.InstanceSource = &source{}

func newServiceStore(c *Sidecar, service models.Service) *source {
	return &source{
		s: service,
		c: c,
	}
}

func (d *source) Instances() ([]models.Instance, error) {
	return d.c.instances(d.s)
}
