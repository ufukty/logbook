package reqs

import (
	"fmt"
	"testing"
)

func TestSetString(t *testing.T) {
	type UserId string
	var dst Cookie[UserId]
	v := "00000000-0000-0000-0000-0000000000000000"

	err := dst.setCookieValue(v)
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}

	if string(dst.Value) != v {
		t.Errorf("assertion, want %q got %q", v, string(dst.Value))
	}
}

func TestSetInt(t *testing.T) {
	type UserId int
	var dst Cookie[UserId]
	v := "0"

	err := dst.setCookieValue(v)
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}

	if fmt.Sprint(dst.Value) != v {
		t.Errorf("assertion, want %q got %q", v, fmt.Sprint(dst.Value))
	}
}
