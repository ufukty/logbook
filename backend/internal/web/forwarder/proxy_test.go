package forwarder

import (
	"fmt"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/models"
	"sync"
	"testing"
	"time"
)

type mockInstanceSource []models.Instance

func (mis *mockInstanceSource) Instances() ([]models.Instance, error) {
	return *mis, nil
}

func TestConcurrentMapWrites(t *testing.T) {
	l, deplcfg, _, err := startup.TestDependencies()
	if err != nil {
		t.Fatal(fmt.Errorf("startup: %w", err))
	}

	is := &mockInstanceSource{}
	for i := 0; i < 100; i++ {
		*is = append(*is, models.Instance{Address: "127.0.0.1", Port: -1})
	}

	proxy := New(is, deplcfg, l)

	var wg sync.WaitGroup
	concurrentGoroutines := 100

	for i := 0; i < concurrentGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_, err := proxy.next()
				if err != nil && err != balancer.ErrNoHostAvailable {
					t.Errorf("Unexpected error: %v", err)
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}

	wg.Wait()
}
