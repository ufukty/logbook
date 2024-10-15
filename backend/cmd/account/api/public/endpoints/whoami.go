package endpoints

import (
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/registration/decls"
	"logbook/models/columns"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type WhoAmIResponse struct {
	Uid       columns.UserId   `json:"uid"`
	Firstname string           `json:"firstname"`
	Lastname  string           `json:"lastname"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (e Endpoints) WhoAmI(id decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("no session_token found")
	}

	sessionToken := cookie.Value

	profile, err := e.a.WhoAmI(r.Context(), columns.SessionToken(sessionToken))
	if err != nil {
		switch err {
		case
			app.ErrProfileNotFound,
			app.ErrUserNotFound,
			app.ErrSessionNotFound:
			http.Error(w, redact(err), http.StatusUnauthorized)
		default:
			http.Error(w, redact(err), http.StatusInternalServerError)
		}
		return fmt.Errorf("app.WhoAmI: %w", err)
	}

	bs := WhoAmIResponse{
		Uid:       profile.Uid,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		CreatedAt: profile.CreatedAt,
	}

	requests.WriteJsonResponse(bs, w)

	return nil
}
