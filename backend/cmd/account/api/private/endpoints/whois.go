package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration/middlewares"
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
func (e Endpoints) WhoIs(rid receptionist.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &WhoIsRequest{}

	err := requests.ParseRequest(w, r, bq)
	if err != nil {
		http.Error(w, fmt.Errorf("ParseRequest :%w", err).Error(), http.StatusInternalServerError)
		return fmt.Errorf("ParseRequest: %w", err)
	}

	profile, err := e.a.WhoIs(r.Context(), bq.SessionToken)
	if err != nil {
		http.Error(w, fmt.Errorf("app.WhoIs :%w", err).Error(), http.StatusInternalServerError)
		return fmt.Errorf("WhoIs: %w", err)
	}

	bs := WhoIsResponse{
		Uid:       profile.Uid,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		CreatedAt: profile.CreatedAt,
	}

	err = requests.WriteJsonResponse(bs, w)
	if err != nil {
		return fmt.Errorf("write json response: %w", err)
	}

	return nil
}
