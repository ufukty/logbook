package database

import "fmt"

type ObjectiveType string

const (
	Rock    = ObjectiveType("rock")
	Regular = ObjectiveType("regular")
)

func (s *ObjectiveType) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("scan source is not []byte")
	}
	switch string(bytes) {
	case "rock":
		*s = Rock
	case "regular":
		*s = Regular
	default:
		return fmt.Errorf("scan source is unknown value %q", bytes)
	}
	return nil
}
