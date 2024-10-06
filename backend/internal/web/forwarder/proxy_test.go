package forwarder

import (
	"logbook/internal/logger"
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
	is := &mockInstanceSource{}
	for i := 0; i < 100; i++ {
		*is = append(*is, models.Instance{Address: "127.0.0.1", Port: -1})
	}
	service := models.Service("test-service")
	servicepath := "/test"

	proxy := New(is, service, servicepath, logger.New("test"))

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
