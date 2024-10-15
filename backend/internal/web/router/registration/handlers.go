package registration

import (
	"fmt"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration/middlewares"
	"net/http"
)

func pong[T any](rid receptionist.RequestId, store *T, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "pong")
	return nil
}

type stripper struct {
	builtin http.Handler
}

func newStripper(prefix string, h http.Handler) *stripper {
	return &stripper{
		builtin: http.StripPrefix(prefix, h),
	}
}

func (s stripper) Strip(rid receptionist.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	s.builtin.ServeHTTP(w, r)
	return nil
}
