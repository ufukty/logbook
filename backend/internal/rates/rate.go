package rates

import (
	"time"
)

func inTheLast(timestamps []time.Time, d time.Duration) int {
	cutoff := time.Now().Add(-1 * d)
	events := 0
	for _, t := range timestamps {
		if t.After(cutoff) {
			events++
		}
	}
	return events
}

type Rates map[time.Duration]int

func IsTheTime(timestamps []time.Time, at Rates) bool {
	for duration, times := range at {
		if inTheLast(timestamps, duration) >= times {
			return false
		}
	}
	return true
}
