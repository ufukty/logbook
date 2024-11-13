package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type ReattachObjectiveRequest struct {
	Subject       columns.ObjectiveId `json:"subject"`
	CurrentParent models.Ovid         `json:"current-parent"`
	NextParent    models.Ovid         `json:"new-parent"`
}

type ReattachObjectiveResponse struct {
	Subject models.Ovid `json:"subject"`
}

// PATCH
func (p *Public) ReattachObjective(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	wi, err := p.sessions.WhoIs(&endpoints.WhoIsRequest{SessionToken: st})
	if err != nil {
		p.l.Println(fmt.Errorf("sessions.WhoIs: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &ReattachObjectiveRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := validate.RequestFields(bq); err != nil {
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = p.a.Reattach(r.Context(), app.ReattachParams{
		Actor:         wi.Uid,
		CurrentParent: bq.CurrentParent,
		NextParent:    bq.NextParent,
		Subject:       bq.Subject,
	})

	bs := ReattachObjectiveResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		return
	}
}
