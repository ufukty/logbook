package endpoints

import (
	"fmt"
	"logbook/cmd/sessions/app"
	"logbook/internal/cookies"
	"logbook/models/columns"
	"net/http"
)

type WhoAmIResponse struct {
	Uid columns.UserId `json:"uid"`
}

// GET
func (p *Public) WhoAmI(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	uid, err := p.a.WhoAmI(r.Context(), st)
	if err != nil {
		p.l.Println(fmt.Errorf("app.WhoAmI: %w", err))
		switch err {
		case
			app.ErrProfileNotFound,
			app.ErrUserNotFound,
			app.ErrSessionNotFound:
			http.Error(w, redact(err), http.StatusUnauthorized)
		default:
			http.Error(w, redact(err), http.StatusInternalServerError)
		}
		return

	}

	bs := WhoAmIResponse{
		Uid: uid,
	}

	bs.Write(w)
}
