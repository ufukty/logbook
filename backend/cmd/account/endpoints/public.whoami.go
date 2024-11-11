package endpoints

import (
	"fmt"
	"logbook/cmd/account/app"
	"logbook/internal/web/requests"
	"logbook/models/columns"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type WhoAmIRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
}

type WhoAmIResponse struct {
	Uid       columns.UserId   `json:"uid"`
	Firstname string           `json:"firstname"`
	Lastname  string           `json:"lastname"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// GET
func (p *Public) WhoAmI(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		p.l.Println(fmt.Errorf("no session_token found"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	profile, err := p.a.WhoAmI(r.Context(), columns.SessionToken(sessionToken))
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
		Uid:       profile.Uid,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		CreatedAt: profile.CreatedAt,
	}

	bs.Write(w)
}
