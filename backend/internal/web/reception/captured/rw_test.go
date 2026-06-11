package captured

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fake
type hijacker struct {
	http.ResponseWriter
}

func (h hijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	panic("not to call")
}

var _ http.Hijacker = (*hijacker)(nil)

func ExampleResponseWriter_additionalMethods() {
	var rw http.ResponseWriter = New(hijacker{})
	if _, ok := rw.(http.Hijacker); !ok {
		fmt.Println("the non-ResponseWriter methods are not available directly")
	}
	if _, ok := rw.(interface{ Unwrap() http.ResponseWriter }).Unwrap().(http.Hijacker); ok {
		fmt.Println("become so, after unwrapping")
	}
	// Output:
	// the non-ResponseWriter methods are not available directly
	// become so, after unwrapping
}

func TestBytes(t *testing.T) {
	tcs := map[uint]string{
		0:                 "0B",
		10:                "10B",
		20:                "20B",
		1000:              "1000B",
		1024:              "1KB",
		1024 + 102:        "1.1KB",
		1000 * 1024:       "1000KB",
		1024 * 1024:       "1MB",
		1.5 * 1024 * 1024: "1.5MB",
		1<<63 - 1:         "8EB",
	}
	for input, expected := range tcs {
		t.Run(fmt.Sprintf("%d", input), func(t *testing.T) {
			got := bytes(input)
			if got != expected {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestResponseWriter_SizeRepresentation(t *testing.T) {
	crw := New(httptest.NewRecorder())
	crw.Write([]byte("lorem ipsum dolor sit amet consectetur adipscing elit"))
	expected, got := "53B", crw.SizeRepresentation()
	if expected != got {
		t.Errorf("expected %q got %q", expected, got)
	}
}

func TestResponseWriter_statusCodeRepresentation(t *testing.T) {
	t.Run("untouched", func(t *testing.T) {
		w := New(httptest.NewRecorder())
		expected, got := "0", w.StatusRepresentation()
		if expected != got {
			t.Errorf("expected %q got %q", expected, got)
		}
	})
	t.Run("implicit", func(t *testing.T) {
		w := New(httptest.NewRecorder())
		w.Write([]byte("lorem ipsum dolor sit amet."))
		expected, got := "200*", w.StatusRepresentation()
		if expected != got {
			t.Errorf("expected %q got %q", expected, got)
		}
	})
	t.Run("explicit", func(t *testing.T) {
		w := New(httptest.NewRecorder())
		w.WriteHeader(200)
		expected, got := "200", w.StatusRepresentation()
		if expected != got {
			t.Errorf("expected %q got %q", expected, got)
		}
	})
}
