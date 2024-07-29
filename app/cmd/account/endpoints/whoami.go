package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/internal/web/requests"
	database "logbook/models/columns"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type WhoAmIResponse struct {
	Uid       database.UserId  `json:"uid"`
	Firstname string           `json:"firstname"`
	Lastname  string           `json:"lastname"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (e Endpoints) WhoAmI(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	profile, err := e.app.WhoAmI(r.Context(), database.SessionToken(sessionToken))
	if err != nil {
		log.Println(fmt.Errorf("app.WhoAmI: %w", err))
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

	requests.WriteJsonResponse(bs, w)
}
