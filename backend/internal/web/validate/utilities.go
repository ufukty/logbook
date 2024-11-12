package validate

import (
	"fmt"
	"regexp"
)

var (
	ErrPattern = fmt.Errorf("content")
	ErrLong    = fmt.Errorf("too long")
	ErrShort   = fmt.Errorf("too short")
	ErrEmpty   = fmt.Errorf("empty")
)

// for external use
func StringBasics(s string, min, max int, pattern *regexp.Regexp) error {
	if l := len(s); l == 0 {
		return ErrEmpty
	} else if l < min {
		return ErrShort
	} else if l > max {
		return ErrLong
	} else if pattern != nil && !pattern.MatchString(string(s)) {
		return ErrPattern
	}
	return nil
}
