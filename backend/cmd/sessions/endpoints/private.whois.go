package endpoints

import (
	"fmt"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"net/http"
)

type WhoIsRequest struct {
	SessionToken columns.SessionToken `json:"session_token"` // body because a service is asking
}

type WhoIsResponse struct {
	Uid columns.UserId `json:"uid"`
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

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	uid, err := p.a.WhoAmI(r.Context(), bq.SessionToken)
	if err != nil {
		p.l.Println(fmt.Errorf("WhoAmI: %w", err))
		http.Error(w, fmt.Errorf("app.WhoIs :%w", err).Error(), http.StatusInternalServerError)
		return
	}

	bs := WhoIsResponse{
		Uid: uid,
	}

	err = bs.Write(w)
	if err != nil {
		p.l.Println(fmt.Errorf("write json response: %w", err))
		return
	}
}
