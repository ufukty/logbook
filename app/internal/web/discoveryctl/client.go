package discoveryctl

import (
	"context"
	"fmt"
	"logbook/cmd/registry/app"
	registry "logbook/cmd/registry/client"
	"logbook/cmd/registry/endpoints"
	"logbook/internal/web/logger"
	"logbook/models"
	"sync"
	"time"
)

// summary:
//   - Pass the list of all services that will be needed through runtime into [New]
//   - [Client] will periodically fetch the list of instances of services and cache them
//   - [Client] will periodically recheck the instance with registry after [Client.SetInstanceDetails] called
//   - [Client] uses [registry.Client] to periodically query the registry service,
//   - [Client.InstanceSource] returns a struct which conforms [balancer.InstanceSource]
type Client struct {
	ctl      *registry.Client
	reload   time.Duration
	services []models.Service

	l      logger.Logger
	ctx    context.Context
	cancel context.CancelFunc

	store   map[models.Service][]models.Instance
	storemu sync.RWMutex

	iid   app.InstanceId
	iidmu sync.RWMutex
}

func New(ctl *registry.Client, period time.Duration, services []models.Service) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Client{
		ctl:      ctl,
		store:    map[models.Service][]models.Instance{},
		services: services,

		l:      *logger.NewLogger("Service Discovery Client"),
		reload: period,
		ctx:    ctx,
		cancel: cancel,
	}
	go d.tick()
	return d
}

func (d *Client) Stop() {
	d.cancel()
}

func (d *Client) queryserver() error {
	for _, service := range d.services {
		bs, err := d.ctl.ListInstances(&endpoints.ListInstancesRequest{Service: service})
		if err != nil {
			return fmt.Errorf("sending listing request: %w", err)
		}
		d.store[service] = bs.Instances
	}
	return nil
}

func (d *Client) recheck() error {
	if d.iid == app.InstanceId("") {
		return nil
	}
	r, err := d.ctl.RecheckInstance(&endpoints.RecheckInstanceRequest{
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

func (d *Client) tick() {
	t := time.NewTicker(d.reload)
	defer t.Stop()
	for {
		select {
		case <-t.C:
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
		case <-d.ctx.Done():
			return
		}
	}
}

func (d *Client) instances(service models.Service) ([]models.Instance, error) {
	d.storemu.RLock()
	defer d.storemu.RUnlock()
	return d.store[service], nil
}

func (c *Client) InstanceSource(s models.Service) *source {
	return newServiceStore(c, s)
}

func (c *Client) SetInstanceDetails(s models.Service, i models.Instance) error {
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
