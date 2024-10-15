package endpoints

import (
	"fmt"
	"logbook/internal/web/router/reception"
	"logbook/models/columns"
	"net/http"
)

// TODO: add anti-CSRF token checks
func (e Endpoints) Logout(id reception.RequestId, store *reception.Store, w http.ResponseWriter, r *http.Request) error {
	st := columns.SessionToken(r.Header.Get("session_token"))

	if err := st.Validate(); err != nil {
		http.Error(w, redact(err), http.StatusUnauthorized)
		return fmt.Errorf("invalid session_token: %w", err)
	}

	err := e.a.Logout(r.Context(), st)
	if err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("saving session deletion to database: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusOK)

	return nil
}
