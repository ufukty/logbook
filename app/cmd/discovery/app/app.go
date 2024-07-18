package app

import (
	"context"
	"fmt"
	"log"
	"logbook/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InstanceId string

type App struct {
	instanceTimeout time.Duration
	clearancePeriod time.Duration

	instances map[models.Service]*Set[InstanceId]
	details   map[InstanceId]models.Instance
	checks    map[InstanceId]time.Time
	cache     map[models.Service][]models.Instance

	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func New(instanceTimeout, clearancePeriod time.Duration) *App {
	ctx, cancel := context.WithCancel(context.Background())
	a := &App{
		instanceTimeout: instanceTimeout,
		clearancePeriod: clearancePeriod,
		instances:       map[models.Service]*Set[InstanceId]{},
		details:         map[InstanceId]models.Instance{},
		checks:          map[InstanceId]time.Time{},
		cache:           map[models.Service][]models.Instance{},
		ctx:             ctx,
		cancel:          cancel,
	}
	go a.ticker()
	return a
}

func (a *App) Stop() {
	a.cancel()
}

func (a *App) RegisterInstance(s models.Service, i models.Instance) (InstanceId, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	t := time.Now()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}
	iid := InstanceId(uuid.String())
	if _, ok := a.instances[s]; !ok {
		a.instances[s] = NewSet[InstanceId]()
	}
	a.instances[s].Add(iid)
	a.checks[iid] = t
	a.details[iid] = i
	a.buildCache(s)
	return iid, nil
}

func (a *App) RecheckInstance(iid InstanceId) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.details[iid]; !ok {
		return fmt.Errorf("instance is either deleted or never registered")
	}
	a.checks[iid] = time.Now()
	return nil
}

func (a *App) ListInstances(s models.Service) ([]models.Instance, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if _, ok := a.instances[s]; !ok {
		return []models.Instance{}, nil
	}
	if list, ok := a.cache[s]; ok {
		return list, nil
	}
	return a.cache[s], nil
}

func (a *App) ticker() {
	ticker := time.NewTicker(a.clearancePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.clearOutdated()
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *App) clearOutdated() {
	a.mu.Lock()
	defer a.mu.Unlock()

	t := time.Now()
	for s, instanceSet := range a.instances {
		toClear := []InstanceId{}
		for _, iid := range instanceSet.Items() {
			if t.Sub(a.checks[iid]) > a.instanceTimeout {
				toClear = append(toClear, iid)
			}
		}
		if len(toClear) > 0 {
			for _, iid := range toClear {
				log.Printf("deleted %q (%s:%s) for %q\n", iid, a.details[iid].Address, a.details[iid].Port, s)
				delete(a.details, iid)
				delete(a.checks, iid)
				a.instances[s].Delete(iid)
			}
			a.buildCache(s)
		}
	}
}

func (a *App) buildCache(s models.Service) {
	cache := []models.Instance{}
	if set, ok := a.instances[s]; ok {
		for _, iid := range set.Items() {
			if instance, ok := a.details[iid]; ok {
				cache = append(cache, instance)
			}
		}
	}
	a.cache[s] = cache
}
