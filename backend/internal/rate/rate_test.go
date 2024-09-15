package rate

import (
	"logbook/internal/average"
	"testing"
	"time"
)

func TestIsAllowedPositive(t *testing.T) {
	tests := map[string]struct {
		timestamps []time.Time
		rates      map[time.Duration]int
	}{
		"No requests, should be allowed": {
			timestamps: []time.Time{},
			rates:      map[time.Duration]int{time.Minute: 5, time.Hour: 100, average.Day: 1000, average.Week: 5000},
		},
		"Within minute limit, should be allowed": {
			timestamps: []time.Time{
				time.Now().Add(-30 * time.Second),
			},
			rates: map[time.Duration]int{time.Minute: 5, time.Hour: 100, average.Day: 1000, average.Week: 5000},
		},
		"Within hour limit, should be allowed": {
			timestamps: []time.Time{
				time.Now().Add(-59 * time.Minute),
				time.Now().Add(-30 * time.Minute),
			},
			rates: map[time.Duration]int{time.Minute: 5, time.Hour: 100, average.Day: 1000, average.Week: 5000},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if !IsTheTime(test.timestamps, test.rates) {
				t.Error("expected 'allowed'")
			}
		})
	}
}
func TestIsAllowedNegative(t *testing.T) {
	tests := map[string]struct {
		timestamps []time.Time
		rate       map[time.Duration]int
	}{
		"Exceeds minute limit, should not be allowed": {
			timestamps: []time.Time{
				time.Now().Add(-50 * time.Second),
				time.Now().Add(-40 * time.Second),
				time.Now().Add(-30 * time.Second),
				time.Now().Add(-20 * time.Second),
				time.Now().Add(-10 * time.Second),
			},
			rate: map[time.Duration]int{time.Minute: 4, time.Hour: 100, average.Day: 1000, average.Week: 5000},
		},
		"Exceeds hour limit, should not be allowed": {
			timestamps: []time.Time{
				time.Now().Add(-59 * time.Minute),
				time.Now().Add(-58 * time.Minute),
				time.Now().Add(-30 * time.Minute),
				time.Now().Add(-15 * time.Minute),
				time.Now().Add(-10 * time.Minute),
				time.Now().Add(-5 * time.Minute),
				time.Now().Add(-1 * time.Minute),
			},
			rate: map[time.Duration]int{time.Minute: 5, time.Hour: 6, average.Day: 1000, average.Week: 5000},
		},
		"Exceeds day limit, should not be allowed": {
			timestamps: []time.Time{
				time.Now().Add(-23 * time.Hour),
				time.Now().Add(-22 * time.Hour),
				time.Now().Add(-21 * time.Hour),
				time.Now().Add(-20 * time.Hour),
				time.Now().Add(-19 * time.Hour),
				time.Now().Add(-18 * time.Hour),
			},
			rate: map[time.Duration]int{time.Minute: 5, time.Hour: 100, average.Day: 5, average.Week: 5000},
		},
		"Exceeds week limit, should not be allowed": {
			timestamps: []time.Time{
				time.Now().Add(-6 * 24 * time.Hour),
				time.Now().Add(-5 * 24 * time.Hour),
				time.Now().Add(-4 * 24 * time.Hour),
				time.Now().Add(-3 * 24 * time.Hour),
				time.Now().Add(-2 * 24 * time.Hour),
				time.Now().Add(-1 * 24 * time.Hour),
			},
			rate: map[time.Duration]int{time.Minute: 5, time.Hour: 100, average.Day: 1000, average.Week: 5},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if IsTheTime(test.timestamps, test.rate) {
				t.Error("expected 'unallowed'")
			}
		})
	}
}
