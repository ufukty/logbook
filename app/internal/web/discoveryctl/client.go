package discoveryctl

import (
	"context"
	"fmt"
	servicereg "logbook/cmd/registry/client"
	"logbook/cmd/registry/endpoints"
	"logbook/internal/web/logger"
	"logbook/models"
	"sync"
	"time"
)

// summary:
//   - instantiate an instance of [Client] per-service to be reached
//   - uses [servicereg.Client] to periodically query the service registry,
//   - contains a cache for the instances,
//   - complies the [balancer.InstanceSource] interface
type Client struct {
	ctl      *servicereg.Client
	store    map[models.Service][]models.Instance
	services []models.Service

	l      logger.Logger
	reload time.Duration
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctl *servicereg.Client, services []models.Service) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Client{
		ctl:      ctl,
		store:    map[models.Service][]models.Instance{},
		services: services,

		l:      *logger.NewLogger("Service Discovery Client"),
		reload: time.Second,
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
		d.store[service] = *bs
	}
	return nil
}

func (d *Client) tick() {
	t := time.NewTicker(d.reload)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			d.mu.Lock()
			if err := d.queryserver(); err != nil {
				d.l.Println(fmt.Errorf("error: querying registry service: %w", err))
			}
			d.mu.Unlock()
		case <-d.ctx.Done():
			return
		}
	}
}

func (d *Client) instances(service models.Service) ([]models.Instance, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.store[service], nil
}

func (c *Client) InstanceSource(s models.Service) *source {
	return newServiceStore(c, s)
}
