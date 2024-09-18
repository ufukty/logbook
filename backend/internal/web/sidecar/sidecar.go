package sidecar

import (
	"context"
	"fmt"
	"logbook/cmd/registry/app"
	registry "logbook/cmd/registry/client"
	"logbook/cmd/registry/endpoints"
	"logbook/config/deployment"
	"logbook/internal/web/logger"
	"logbook/models"
	"sync"
	"time"
)

// summary:
//   - Pass the list of all services that will be needed through runtime into [New]
//   - [Sidecar] will periodically fetch the list of instances of services and cache them
//   - [Sidecar] will periodically recheck the instance with registry after [Sidecar.SetInstanceDetails] called
//   - [Sidecar] uses [registry.Client] to periodically query the registry service,
//   - [Sidecar.InstanceSource] returns a struct which conforms [balancer.InstanceSource]
type Sidecar struct {
	ctl      *registry.Client
	deplycfg *deployment.Config

	service  models.Service
	services []models.Service

	l      logger.Logger
	ctx    context.Context
	cancel context.CancelFunc

	store   map[models.Service][]models.Instance
	storemu sync.RWMutex

	iid   app.InstanceId
	iidmu sync.RWMutex
}

func New(ctl *registry.Client, deplycfg *deployment.Config, services []models.Service) *Sidecar {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Sidecar{
		ctl:      ctl,
		store:    map[models.Service][]models.Instance{},
		services: services,

		l:        *logger.NewLogger("Sidecar"),
		deplycfg: deplycfg,
		ctx:      ctx,
		cancel:   cancel,
	}
	go d.tick()
	return d
}

func (d *Sidecar) Stop() {
	d.cancel()
}

func (d *Sidecar) queryserver() error {
	for _, service := range d.services {
		d.l.Printf("queryserver for %s\n", service)
		bs, err := d.ctl.ListInstances(&endpoints.ListInstancesRequest{Service: service})
		if err != nil {
			return fmt.Errorf("sending listing request: %w", err)
		}
		d.store[service] = bs.Instances
	}
	return nil
}

func (d *Sidecar) recheck() error {
	if d.iid == app.InstanceId("") || d.service == "" { // sidecar without registration (eg. "api-gateway")
		return nil
	}
	d.l.Println("rechecking...")
	r, err := d.ctl.RecheckInstance(&endpoints.RecheckInstanceRequest{
		Service:    d.service,
		InstanceId: d.iid,
	})
	if err != nil {
		return fmt.Errorf("registry.Client.RecheckInstance: %w", err)
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("registry service returned non-200 status code: %d", r.StatusCode)
	}
	return nil
}

func (d *Sidecar) update() {
	d.storemu.Lock()
	if err := d.queryserver(); err != nil {
		d.l.Println(fmt.Errorf("tick: queryserver: %w", err))
	}
	d.storemu.Unlock()

	d.iidmu.RLock()
	if err := d.recheck(); err != nil {
		d.l.Println(fmt.Errorf("tick: recheck: %w", err))
	}
	d.iidmu.RUnlock()
}

func (d *Sidecar) tick() {
	time.Sleep(d.deplycfg.Sidecar.TickerDelay)
	t := time.NewTicker(d.deplycfg.Sidecar.TickerPeriod)
	d.update() // before the first tick
	defer t.Stop()
	for {
		select {
		case <-t.C:
			d.update()
		case <-d.ctx.Done():
			return
		}
	}
}

func (d *Sidecar) instances(service models.Service) ([]models.Instance, error) {
	d.storemu.RLock()
	defer d.storemu.RUnlock()
	return d.store[service], nil
}

func (c *Sidecar) InstanceSource(s models.Service) *source {
	return newServiceStore(c, s)
}

func (c *Sidecar) SetInstanceDetails(s models.Service, i models.Instance) error {
	c.service = s
	c.l.Printf("registering the instance: %s -> %s\n", s, i)
	r, err := c.ctl.RegisterInstance(&endpoints.RegisterInstanceRequest{
		Service: s,
		TLS:     i.Tls,
		Address: i.Address,
		Port:    i.Port,
	})
	if err != nil {
		return fmt.Errorf("registry.Client.RegisterInstance: %w", err)
	}
	c.iidmu.Lock()
	defer c.iidmu.Unlock()
	c.iid = r.InstanceId
	return nil
}
