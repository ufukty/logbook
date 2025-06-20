package endpoints

import (
	"fmt"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

// TODO:
type TagAssignRequest struct {
	Subject models.Ovid   `json:"subject"`
	Tid     columns.TagId `json:"tid"`
}

// POST
func (p *Public) TagAssign(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	_, err = p.sessions.WhoIs(&endpoints.WhoIsRequest{
		SessionToken: st,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("sessions.WhoIs: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &TagAssignRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	panic("to implement") // TODO:
}
