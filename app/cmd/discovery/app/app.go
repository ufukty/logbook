package app

import (
	"fmt"
	"log"
	"logbook/models"
	"time"

	"github.com/google/uuid"
)

type Instance struct {
	TLS     bool
	Address string
	Port    string
}

type InstanceId string

type App struct {
	instances        map[models.Service]*Set[InstanceId]
	details          map[InstanceId]Instance
	checks           map[InstanceId]time.Time
	housekeepingtime map[models.Service]time.Time
	housekeepinglock map[models.Service]bool

	cache map[models.Service][]Instance
}

func New() *App {
	return &App{
		instances:        map[models.Service]*Set[InstanceId]{},
		details:          map[InstanceId]Instance{},
		checks:           map[InstanceId]time.Time{},
		housekeepingtime: map[models.Service]time.Time{},
		housekeepinglock: map[models.Service]bool{},
		cache:            map[models.Service][]Instance{},
	}
}

func (a *App) RegisterInstance(s models.Service, i Instance) (InstanceId, error) {
	t := time.Now()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generating uuid: %w", err)
	}
	iid := InstanceId(uuid.String())
	if _, ok := a.instances[s]; !ok {
		a.instances[s] = &Set[InstanceId]{}
	}
	a.instances[s].Add(iid)
	a.checks[iid] = t
	a.details[iid] = i
	return iid, nil
}

func (a *App) RecheckInstance(iid InstanceId) error {
	if _, ok := a.details[iid]; !ok {
		return fmt.Errorf("instance is either deleted or never registered")
	}
	a.checks[iid] = time.Now()
	return nil
}

func (a *App) clearOutdated(s models.Service) {
	if _, ok := a.housekeepinglock[s]; ok {
		return
	}
	defer delete(a.housekeepinglock, s)

	t := time.Now()
	if l, ok := a.housekeepingtime[s]; ok && l.Sub(t) < time.Minute { // house keeping period
		return
	}
	if _, ok := a.instances[s]; !ok {
		return
	}
	toClear := []InstanceId{}
	for iid := range *a.instances[s] {
		if t.Sub(a.checks[iid]) < time.Minute { //
			toClear = append(toClear, iid)
		}
	}
	if len(toClear) > 0 {
		delete(a.cache, s)
		for _, iid := range toClear {
			log.Printf("deleted %q (%s:%s) for %q\n", iid, a.details[iid].Address, a.details[iid].Port, s)
			delete(a.details, iid)
			delete(a.checks, iid)
			a.instances[s].Delete(iid)
		}
	}
	a.housekeepingtime[s] = time.Now()

	return
}

func (a *App) buildCache(s models.Service) {
	cache := []Instance{}
	if _, ok := a.instances[s]; !ok {
		a.cache[s] = cache
		return
	}
	for iid := range *a.instances[s] {
		if i, ok := a.details[iid]; ok {
			cache = append(cache, i)
		}
	}
	a.cache[s] = cache
}

func (a *App) ListInstances(s models.Service) ([]Instance, error) {
	if _, ok := a.instances[s]; !ok {
		return []Instance{}, nil
	}
	a.clearOutdated(s)
	if list, ok := a.cache[s]; ok {
		return list, nil
	}
	a.buildCache(s)
	return a.cache[s], nil
}
