package discoveryctl

import (
	"context"
	"fmt"
	servicereg "logbook/cmd/registry/client"
	"logbook/cmd/registry/endpoints"
	"logbook/internal/web/balancer"
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
	ctl     *servicereg.Client
	store   []models.Instance
	service models.Service

	l      logger.Logger
	reload time.Duration
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

var _ balancer.InstanceSource = &Client{}

func New(ctl *servicereg.Client, service models.Service) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Client{
		ctl:     ctl,
		store:   []models.Instance{},
		service: service,

		l:      *logger.NewLogger("Discover"),
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
	bs, err := d.ctl.ListInstances(&endpoints.ListInstancesRequest{Service: d.service})
	if err != nil {
		return fmt.Errorf("sending listing request: %w", err)
	}
	d.store = *bs
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
				d.l.Println(fmt.Errorf("error: querying service discovery service: %w", err))
			}
			d.mu.Unlock()
		case <-d.ctx.Done():
			return
		}
	}
}

func (d *Client) Instances() ([]models.Instance, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.store, nil
}
