package validate

import "time"

func TimeBasics(t, min, max time.Time) bool {
	return (t.After(min) || t.Equal(min)) && (t.Before(max) || t.Equal(max))
}
