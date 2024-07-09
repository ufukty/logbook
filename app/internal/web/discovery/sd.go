package discovery

import (
	"fmt"
	"logbook/internal/web/discovery/models/local"
	"logbook/internal/web/discovery/models/stage"
	"logbook/models"
	"sync"
	"time"
)

type Pool interface {
	ServicePool(s models.Service) ([]string, error)
}

type ConfigBasedServiceDiscovery struct {
	e            models.Environment
	pool         Pool
	configPath   string
	updateLock   sync.Mutex
	updatePeriod time.Duration
}

func New(e models.Environment, configPath string, updatePeriod time.Duration) *ConfigBasedServiceDiscovery {
	sd := ConfigBasedServiceDiscovery{
		e:            e,
		configPath:   configPath,
		updatePeriod: updatePeriod,
	}
	sd.readConfig()
	go sd.tick()
	return &sd
}

func (sd *ConfigBasedServiceDiscovery) readConfig() {
	if !sd.updateLock.TryLock() {
		return
	}
	var err error
	switch sd.e {
	case models.Local:
		sd.pool, err = local.ReadLocal(sd.configPath)
		if err != nil {
			panic(fmt.Errorf("reading service discovery file for local environment: %w", err))
		}
	case models.Stage:
		sd.pool, err = stage.ReadStage(sd.configPath)
		if err != nil {
			panic(fmt.Errorf("reading service discovery file for stage environment: %w", err))
		}
	}
	sd.updateLock.Unlock()
}

func (sd *ConfigBasedServiceDiscovery) tick() {
	for range time.Tick(sd.updatePeriod) {
		sd.readConfig()
	}
}

func (sd *ConfigBasedServiceDiscovery) ServicePool(service models.Service) ([]string, error) {
	return sd.pool.ServicePool(service)
}
