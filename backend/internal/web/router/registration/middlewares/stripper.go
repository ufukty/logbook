package middlewares

import (
	"logbook/internal/web/router/registration/receptionist/decls"
	"net/http"
)

type stripper struct {
	builtin http.Handler
}

func NewStripper(prefix string, h http.Handler) *stripper {
	return &stripper{
		builtin: http.StripPrefix(prefix, h),
	}
}

func (s stripper) Strip(rid decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	s.builtin.ServeHTTP(w, r)
	return nil
}
