package challenge

import (
	"fmt"
	"strconv"
	"testing"
)

func TestRandString(t *testing.T) {
	for i := range 10 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s, err := randstring(alphabet, i)
			if err != nil {
				t.Fatal(fmt.Errorf("act: %w", err))
			}
			t.Log(s)
			if len(s) != i {
				t.Fatalf("assert len, expected %d, got %d", i, len(s))
			}
		})
	}
}
