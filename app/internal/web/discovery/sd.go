package discovery

import (
	"fmt"
	"logbook/models"
	"sync"
	"time"
)

type Config interface {
	ServicePool(s models.Service) ([]string, error)
}

type ServiceDiscovery struct {
	config       Config `yaml:",inline"`
	configPath   string
	updateLock   sync.Mutex
	updatePeriod time.Duration
}

func New(configPath string, updatePeriod time.Duration) *ServiceDiscovery {
	sd := ServiceDiscovery{
		configPath:   configPath,
		updatePeriod: updatePeriod,
	}
	sd.readConfig()
	go sd.tick()
	return &sd
}

func (sd *ServiceDiscovery) readConfig() {
	if !sd.updateLock.TryLock() {
		return
	}
	var err error
	sd.config, err = ReadStage(sd.configPath)
	if err != nil {
		panic(fmt.Errorf("reading service discovery config file: %w", err))
	}
	sd.updateLock.Unlock()
}

func (sd *ServiceDiscovery) tick() {
	for range time.Tick(sd.updatePeriod) {
		sd.readConfig()
	}
}

func (sd *ServiceDiscovery) ServicePool(service models.Service) ([]string, error) {
	return sd.config.ServicePool(service)
}
