package reception

import (
	"fmt"
	"net/http"
	"strings"
)

type cors struct {
	origin  string
	methods string
	headers string
	next    http.HandlerFunc
}

func newCors(next http.HandlerFunc, origin string, methods, headers []string) *cors {
	methods = append(methods, "OPTIONS")
	return &cors{
		origin:  origin,
		methods: strings.Join(methods, ", "),
		headers: strings.Join(headers, ", "),
		next:    next,
	}
}

func (c *cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", c.origin)
	w.Header().Set("Access-Control-Allow-Methods", c.methods)
	w.Header().Set("Access-Control-Allow-Headers", c.headers)

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	c.next(w, r)
}

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
	w.WriteHeader(http.StatusOK)
}
