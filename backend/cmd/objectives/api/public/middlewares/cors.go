package middlewares

import (
	"logbook/internal/web/router/receptionist"
	"net/http"
	"strings"
)

type Cors struct {
	origin  string
	methods string
	headers string
}

func newCors(origin string, methods, headers []string) *Cors {
	methods = append(methods, "OPTIONS")
	return &Cors{
		origin:  origin,
		methods: strings.Join(methods, ", "),
		headers: strings.Join(headers, ", "),
	}
}

func (c *Cors) Handle(id receptionist.RequestId, store *Store, w http.ResponseWriter, r *http.Request) error {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", c.origin)
	w.Header().Set("Access-Control-Allow-Methods", c.methods)
	w.Header().Set("Access-Control-Allow-Headers", c.headers)

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return receptionist.ErrSilent
	}

	return nil
}

type CorsManager struct {
	origin string
}

func NewCorsManager(origin string) *CorsManager {
	return &CorsManager{
		origin: origin,
	}
}

func (c CorsManager) Instantiate(methods, headers []string) *Cors {
	return newCors(c.origin, methods, headers)
}
