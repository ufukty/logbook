package app

import (
	"context"
	"fmt"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/stores"
	"logbook/models"
	"logbook/models/columns"
	"sync"
	"time"
)

type InstanceId string

// maintains per-service locks
type serviceRegistry struct {
	deplycfg *deployment.Config
	l        *logger.Logger

	ctx    context.Context
	cancel context.CancelFunc

	instances *stores.KV[InstanceId, models.Instance]
	checks    *stores.KV[InstanceId, time.Time]
	cache     []models.Instance
	mu        sync.RWMutex
}

func newServiceRegistry(deplycfg *deployment.Config, logname string) *serviceRegistry {
	ctx, cancel := context.WithCancel(context.Background())
	sr := &serviceRegistry{
		deplycfg: deplycfg,
		l:        logger.New(logname),

		ctx:    ctx,
		cancel: cancel,

		instances: stores.NewKV[InstanceId, models.Instance](),
		checks:    stores.NewKV[InstanceId, time.Time](),

		cache: []models.Instance{},
	}
	go sr.ticker()
	return sr
}

func (sr *serviceRegistry) RegisterInstance(i models.Instance) (InstanceId, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	var iid InstanceId
	var err error
	for iid == "" || sr.checks.Has(iid) { // collision
		iid, err = columns.NewUuidV4[InstanceId]()
		if err != nil {
			return "", fmt.Errorf("NewUuidV4[InstanceId]: %w", err)
		}
	}
	sr.l.Printf("welcome: %s (%s)\n", iid, i)
	sr.checks.Set(iid, time.Now())
	sr.instances.Set(iid, i)
	sr.cache = append(sr.cache, i)

	return iid, nil
}

func (sr *serviceRegistry) RecheckInstance(iid InstanceId) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if !sr.checks.Has(iid) {
		return fmt.Errorf("instance is either deleted or never registered")
	}
	sr.l.Printf("welcome back: %s\n", iid)
	sr.checks.Set(iid, time.Now())
	return nil
}

func (sr *serviceRegistry) ListInstances() ([]models.Instance, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	return sr.cache, nil
}

func (sr *serviceRegistry) ticker() {
	ticker := time.NewTicker(sr.deplycfg.Registry.ClearancePeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			sr.houseKeeping()
		case <-sr.ctx.Done():
			return
		}
	}
}

func (sr *serviceRegistry) houseKeeping() {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.l.Println("HouseKeeping starts...")
	t := time.Now()
	toClear := []InstanceId{}
	for iid, last := range sr.checks.Iter() {
		if t.Sub(last) > sr.deplycfg.Registry.InstanceTimeout {
			toClear = append(toClear, iid)
		}
	}
	if len(toClear) > 0 {
		for _, iid := range toClear {
			i, _ := sr.instances.Get(iid)
			sr.l.Printf("instance timeout: %s (%s)\n", iid, i)
			sr.instances.Delete(iid)
			sr.checks.Delete(iid)
		}
		sr.cache = []models.Instance{}
		for _, i := range sr.instances.Iter() {
			sr.cache = append(sr.cache, i)
		}
	}
}

func (sr *serviceRegistry) Stop() {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	sr.cancel()
}

// first contact point before per-service registries
type App struct {
	deplycfg   *deployment.Config
	l          *logger.Logger
	registries *stores.KV[models.Service, *serviceRegistry]
}

func New(deplycfg *deployment.Config) *App {
	a := &App{
		deplycfg:   deplycfg,
		l:          logger.New("Hub"),
		registries: stores.NewKV[models.Service, *serviceRegistry](),
	}
	return a
}

func (a *App) Stop() {
	for _, sr := range a.registries.Iter() {
		sr.Stop()
	}
}

func (a *App) RegisterInstance(s models.Service, i models.Instance) (InstanceId, error) {
	if !a.registries.Has(s) {
		a.l.Println("service registry has generated for:", s)
		a.registries.Set(s, newServiceRegistry(a.deplycfg, fmt.Sprintf("ServiceRegistry(%s)", s)))
	}
	sr, _ := a.registries.Get(s)
	return sr.RegisterInstance(i)
}

func (a *App) RecheckInstance(s models.Service, iid InstanceId) error {
	sr, ok := a.registries.Get(s)
	if !ok {
		return fmt.Errorf("the service %s is not created. is the instance registered itself before try to recheck?", s)
	}
	return sr.RecheckInstance(iid)
}

func (a *App) ListInstances(s models.Service) ([]models.Instance, error) {
	sr, ok := a.registries.Get(s)
	if !ok {
		return nil, fmt.Errorf("the service %s is not created. is the instance registered itself before try to recheck?", s)
	}
	return sr.ListInstances()
}
