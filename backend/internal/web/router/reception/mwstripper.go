package reception

import (
	"net/http"
)

type stripper struct {
	builtin http.Handler
}

func newStripper(prefix string, h http.Handler) *stripper {
	return &stripper{
		builtin: http.StripPrefix(prefix, h),
	}
}

func (s stripper) Strip(rid RequestId, store *Store, w http.ResponseWriter, r *http.Request) error {
	s.builtin.ServeHTTP(w, r)
	return nil
}
