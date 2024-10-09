package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/models/columns"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type WhoIsRequest struct {
	SessionToken columns.SessionToken `json:"session_token"`
}

type WhoIsResponse struct {
	Uid       columns.UserId   `json:"uid"`
	Firstname string           `json:"firstname"`
	Lastname  string           `json:"lastname"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// TODO: What it should return for missing body, invalid token or IO errors?
func (e Endpoints) WhoIs(w http.ResponseWriter, r *http.Request) {
	bq := &WhoIsRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, fmt.Errorf("ParseRequest :%w", err).Error(), http.StatusInternalServerError)
		return
	}

	profile, err := e.a.WhoIs(r.Context(), bq.SessionToken)
	if err != nil {
		http.Error(w, fmt.Errorf("app.WhoIs :%w", err).Error(), http.StatusInternalServerError)
		return
	}

	bs := WhoIsResponse{
		Uid:       profile.Uid,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		CreatedAt: profile.CreatedAt,
	}

	requests.WriteJsonResponse(bs, w)
}
