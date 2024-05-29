package servdisc

import (
	"fmt"
	"logbook/internal/web/servdisc/models/stage"
	"sync"
	"time"
)

type Config interface {
	ServicePool(service string) ([]string, error)
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
	sd.config, err = stage.ReadConfig(sd.configPath)
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

func (sd *ServiceDiscovery) ServicePool(service string) ([]string, error) {
	return sd.config.ServicePool(service)
}
