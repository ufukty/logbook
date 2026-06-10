package captured

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
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
