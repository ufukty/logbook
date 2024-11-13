package endpoints

import (
	"fmt"
	"logbook/internal/cookies"
	"net/http"
)

// TODO: add anti-CSRF token checks
// POST
func (p *Public) Logout(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err = p.a.Logout(r.Context(), st)
	if err != nil {
		p.l.Println(fmt.Errorf("saving session deletion to database: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	cookies.ExpireSessionToken(w)
	w.WriteHeader(http.StatusOK)
}
