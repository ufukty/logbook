package app_test

import (
	"fmt"
	"logbook/cmd/registry/app"
	"logbook/config/deployment"
	"logbook/models"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Test RegisterInstance method
func TestRegisterInstance(t *testing.T) {
	deplcfg, err := deployment.ReadConfig(filepath.Join(os.Getenv("WORKSPACE"), "platform/local/deployment.yml"))
	if err != nil {
		t.Fatal(fmt.Errorf("prep, deployment.ReadConfig: %w", err))
	}
	a := app.New(deplcfg)
	defer a.Stop()

	service := models.Service("test-service")
	instance := models.Instance{ /* Fill with mock fields */ }

	iid, err := a.RegisterInstance(service, instance)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if iid == "" {
		t.Errorf("expected instance ID to be generated, but got empty")
	}
}

// Test RecheckInstance method
func TestRecheckInstance(t *testing.T) {
	deplcfg, err := deployment.ReadConfig(filepath.Join(os.Getenv("WORKSPACE"), "platform/local/deployment.yml"))
	if err != nil {
		t.Fatal(fmt.Errorf("prep, deployment.ReadConfig: %w", err))
	}
	a := app.New(deplcfg)
	defer a.Stop()

	service := models.Service("test-service")
	instance := models.Instance{ /* Fill with mock fields */ }

	iid, err := a.RegisterInstance(service, instance)
	if err != nil {
		t.Errorf("expected no error during instance registration, got %v", err)
	}

	err = a.RecheckInstance(service, iid)
	if err != nil {
		t.Errorf("expected no error during instance recheck, got %v", err)
	}
}

// Test ListInstances method
func TestListInstances(t *testing.T) {
	deplcfg, err := deployment.ReadConfig(filepath.Join(os.Getenv("WORKSPACE"), "platform/local/deployment.yml"))
	if err != nil {
		t.Fatal(fmt.Errorf("prep, deployment.ReadConfig: %w", err))
	}
	a := app.New(deplcfg)
	defer a.Stop()

	service := models.Service("test-service")
	instance1 := models.Instance{ /* Fill with mock fields */ }
	instance2 := models.Instance{ /* Fill with mock fields */ }

	_, err = a.RegisterInstance(service, instance1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = a.RegisterInstance(service, instance2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	instances, err := a.ListInstances(service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(instances) != 2 {
		t.Errorf("expected 2 instances, got %d", len(instances))
	}
}

// Test HouseKeeping logic (clearing timed-out instances)
func TestHouseKeeping(t *testing.T) {
	deplcfg := &deployment.Config{}
	deplcfg.Registry.InstanceTimeout = 50 * time.Millisecond
	deplcfg.Registry.ClearancePeriod = 50 * time.Millisecond
	deplcfg.Registry.ClearanceDelay = 50 * time.Millisecond

	a := app.New(deplcfg)
	defer a.Stop()

	service := models.Service("test-service")
	instance := models.Instance{ /* Fill with mock fields */ }

	_, err := a.RegisterInstance(service, instance)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Sleep for longer than the instance timeout to ensure it is cleared
	time.Sleep(200 * time.Millisecond)

	instances, err := a.ListInstances(service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(instances) != 0 {
		t.Errorf("expected 0 instances after housekeeping, got %d", len(instances))
	}
}
