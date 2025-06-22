package scalars

import (
	"fmt"
	"strconv"
)

func (p *PlacementLength) FromRoute(src string) error {
	a, err := strconv.Atoi(src)
	if err != nil {
		return fmt.Errorf("atoi: %w", err)
	}
	*p = PlacementLength(a)
	return nil
}

func (p *PlacementStart) FromRoute(src string) error {
	a, err := strconv.Atoi(src)
	if err != nil {
		return fmt.Errorf("atoi: %w", err)
	}
	*p = PlacementStart(a)
	return nil
}

func (p PlacementLength) ToRoute() (string, error) { return strconv.Itoa(int(p)), nil }
func (p PlacementStart) ToRoute() (string, error)  { return strconv.Itoa(int(p)), nil }
