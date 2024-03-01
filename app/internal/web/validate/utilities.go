package validate

import (
	"fmt"
	"regexp"
)

var (
	ErrPattern = fmt.Errorf("given value doesn't fit into expected pattern")
	ErrLong    = fmt.Errorf("given value is longer than expected")
	ErrShort   = fmt.Errorf("given value is shorter than expected")
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
