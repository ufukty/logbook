package app

import (
	"logbook/models"
	"testing"
	"time"
)

func TestRegisterInstance(t *testing.T) {
	app := New(100*time.Millisecond, 200*time.Millisecond)
	defer app.Stop()
	service := models.Service("test-service")
	instance := models.Instance{Tls: true, Address: "127.0.0.1", Port: 8080}

	instanceID, err := app.RegisterInstance(service, instance)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if instanceID == "" {
		t.Fatalf("expected non-empty instance ID")
	}
}

func TestRecheckInstance(t *testing.T) {
	app := New(100*time.Millisecond, 200*time.Millisecond)
	defer app.Stop()
	service := models.Service("test-service")
	instance := models.Instance{Tls: true, Address: "127.0.0.1", Port: 8080}

	instanceID, err := app.RegisterInstance(service, instance)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = app.RecheckInstance(instanceID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClearOutdated(t *testing.T) {
	app := New(100*time.Millisecond, 200*time.Millisecond)
	defer app.Stop()
	service := models.Service("test-service")
	instance := models.Instance{Tls: true, Address: "127.0.0.1", Port: 8080}

	// Register an instance and advance time to simulate expiration
	instanceID, err := app.RegisterInstance(service, instance)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Directly manipulate time for testing purposes
	app.mu.Lock()
	app.checks[instanceID] = time.Now().Add(-2 * time.Minute) // Set the last check time to 2 minutes ago
	app.mu.Unlock()

	app.clearOutdated()

	// Check that the instance is removed
	app.mu.RLock()
	defer app.mu.RUnlock()
	_, exists := app.details[instanceID]
	if exists {
		t.Fatalf("expected instance to be removed")
	}
}

func TestBuildCache(t *testing.T) {
	app := New(100*time.Millisecond, 200*time.Millisecond)
	defer app.Stop()
	service := models.Service("test-service")
	instance := models.Instance{Tls: true, Address: "127.0.0.1", Port: 8080}

	_, err := app.RegisterInstance(service, instance)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	app.buildCache(service)

	app.mu.RLock()
	defer app.mu.RUnlock()
	cache, ok := app.cache[service]
	if !ok {
		t.Fatalf("expected cache to be built")
	}
	if len(cache) != 1 {
		t.Fatalf("expected cache length to be 1, got %d", len(cache))
	}
	if cache[0] != instance {
		t.Fatalf("expected cache instance to be %v, got %v", instance, cache[0])
	}
}

func TestListInstances(t *testing.T) {
	app := New(100*time.Millisecond, 200*time.Millisecond)
	defer app.Stop()
	service := models.Service("test-service")
	instance := models.Instance{Tls: true, Address: "127.0.0.1", Port: 8080}

	_, err := app.RegisterInstance(service, instance)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instances, err := app.ListInstances(service)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 1 {
		t.Fatalf("expected instance length to be 1, got %d", len(instances))
	}
	if instances[0] != instance {
		t.Fatalf("expected instance to be %v, got %v", instance, instances[0])
	}
}

func TestTickerClearOutdated(t *testing.T) {
	app := New(100*time.Millisecond, 200*time.Millisecond)
	defer app.Stop()
	service := models.Service("test-service")
	instance := models.Instance{Tls: true, Address: "127.0.0.1", Port: 8080}

	_, err := app.RegisterInstance(service, instance)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Simulate the ticker functionality
	time.Sleep(300 * time.Millisecond) // Wait for the instance to expire

	instances, err := app.ListInstances(service)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 0 {
		t.Fatalf("expected instance length to be 0, got %d", len(instances))
	}
}
