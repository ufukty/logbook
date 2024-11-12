package endpoints

import (
	"fmt"
	"logbook/models/columns"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type WhoIsRequest struct {
	SessionToken columns.SessionToken `json:"session_token"` // body because a service is asking
}

type WhoIsResponse struct {
	Uid       columns.UserId   `json:"uid"`
	Firstname string           `json:"firstname"`
	Lastname  string           `json:"lastname"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// TODO: What it should return for missing body, invalid token or IO errors?
// POST
func (p *Private) WhoIs(w http.ResponseWriter, r *http.Request) {
	bq := &WhoIsRequest{}

	err := bq.Parse(r)
	if err != nil {
		p.l.Println(fmt.Errorf("ParseRequest: %w", err))
		http.Error(w, fmt.Errorf("ParseRequest :%w", err).Error(), http.StatusInternalServerError)
		return
	}

	profile, err := p.a.WhoAmI(r.Context(), bq.SessionToken)
	if err != nil {
		p.l.Println(fmt.Errorf("WhoAmI: %w", err))
		http.Error(w, fmt.Errorf("app.WhoIs :%w", err).Error(), http.StatusInternalServerError)
		return
	}

	bs := WhoIsResponse{
		Uid:       profile.Uid,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		CreatedAt: profile.CreatedAt,
	}

	err = bs.Write(w)
	if err != nil {
		p.l.Println(fmt.Errorf("write json response: %w", err))
		return
	}
}
