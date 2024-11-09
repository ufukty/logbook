package endpoints

import (
	"fmt"
	"logbook/models/columns"
	"net/http"
)

// TODO: add anti-CSRF token checks
func (e Endpoints) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		e.l.Println(fmt.Errorf("no session_token found"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	st := columns.SessionToken(cookie.Value)

	if err := st.Validate(); err != nil {
		e.l.Println(fmt.Errorf("invalid session_token: %w", err))
		http.Error(w, redact(err), http.StatusUnauthorized)
		return
	}

	err = e.a.Logout(r.Context(), st)
	if err != nil {
		e.l.Println(fmt.Errorf("saving session deletion to database: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusOK)
}
