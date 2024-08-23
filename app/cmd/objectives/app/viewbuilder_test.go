package app

import (
	"testing"
)

func TestDoOverlapPositive(t *testing.T) {
	type tc struct {
		a, b int32
		x, y int32
	}
	tcs := map[string]tc{
		"total overlap":    {0, 1, 0, 1},
		"second is nested": {1, 5, 2, 4},
		"first is nested":  {2, 4, 1, 5},
		"mutual left":      {1, 5, 1, 3},
		"mutual right":     {1, 5, 2, 5},
	}
	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			if !doOverlap(tc.a, tc.b, tc.x, tc.y) {
				t.Errorf("returned false for %s", tn)
			}
		})
	}

}

func TestDoOverlapNegative(t *testing.T) {
	type tc struct {
		a, b int32
		x, y int32
	}
	tcs := map[string]tc{
		"second comes after": {1, 5, 6, 7},
		"first comes after":  {6, 7, 1, 5},
		"dimensionless":      {99, 99, 99, 99},
		"dimensionless-2":    {-99, -99, 99, 99},
	}
	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			if doOverlap(tc.a, tc.b, tc.x, tc.y) {
				t.Errorf("returned true for %s", tn)
			}
		})
	}
}
