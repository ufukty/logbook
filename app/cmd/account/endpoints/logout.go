package endpoints

import (
	"fmt"
	"logbook/cmd/account/database"
	"net/http"
)

// TODO: add anti-CSRF token checks
func (e Endpoints) Logout(w http.ResponseWriter, r *http.Request) {
	st := database.SessionToken(r.Header.Get("session_token"))

	if err := st.Validate(); err != nil {
		e.l.Println(fmt.Errorf("invalid session_token: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := e.app.Logout(r.Context(), st)
	if err != nil {
		e.l.Println(fmt.Errorf("saving session deletion to database: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	})
}
