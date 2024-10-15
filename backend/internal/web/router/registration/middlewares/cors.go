package middlewares

import (
	"logbook/internal/web/router/registration/decls"
	"net/http"
	"strings"
)

type cors struct {
	origin  string
	methods string
	headers string
}

func newCors(origin string, methods, headers []string) *cors {
	methods = append(methods, "OPTIONS")
	return &cors{
		origin:  origin,
		methods: strings.Join(methods, ", "),
		headers: strings.Join(headers, ", "),
	}
}

func (c *cors) Handle(id decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", c.origin)
	w.Header().Set("Access-Control-Allow-Methods", c.methods)
	w.Header().Set("Access-Control-Allow-Headers", c.headers)

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return decls.ErrEarlyReturn
	}

	return nil
}

type corsManager struct {
	origin string
}

func NewCorsManager(origin string) *corsManager {
	return &corsManager{
		origin: origin,
	}
}

func (c corsManager) Instantiate(methods, headers []string) *cors {
	return newCors(c.origin, methods, headers)
}
