package discovery

import (
	"context"
	"fmt"
	"logbook/cmd/discovery/endpoints"
	"logbook/internal/web/logger"
	"logbook/models"
	"sync"
	"time"
)

// wrapper for the [Client] which
//   - periodically queries the service-discovery service,
//   - contains cache for the instances,
//   - complies the [balancer.InstanceSource] interface
type DiscoveryStore struct {
	ctl     *Client
	store   []models.Instance
	service models.Service

	l      logger.Logger
	reload time.Duration
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewDiscoveryStore(ctl *Client, service models.Service) *DiscoveryStore {
	ctx, cancel := context.WithCancel(context.Background())
	d := &DiscoveryStore{
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

func (d *DiscoveryStore) Stop() {
	d.cancel()
}

func (d *DiscoveryStore) queryserver() error {
	bs, err := d.ctl.ListInstances(&endpoints.ListInstancesRequest{Service: d.service})
	if err != nil {
		return fmt.Errorf("sending listing request: %w", err)
	}
	d.store = *bs
	return nil
}

func (d *DiscoveryStore) tick() {
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

func (d *DiscoveryStore) Instances() ([]models.Instance, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.store, nil
}
