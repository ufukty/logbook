package validate

import (
	"fmt"
	"regexp"
)

var (
	ErrPattern = fmt.Errorf("content")
	ErrLong    = fmt.Errorf("too long")
	ErrShort   = fmt.Errorf("too short")
)

// for external use
func StringBasics(s string, min, max int, pattern *regexp.Regexp) error {
	if len(s) < min {
		return ErrShort
	} else if len(s) > max {
		return ErrLong
	} else if !pattern.MatchString(string(s)) {
		return ErrPattern
	}
	return nil
}
