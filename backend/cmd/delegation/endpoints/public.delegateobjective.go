package endpoints

import (
	"fmt"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/internal/web/serialize"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type DelegateObjectiveRequest struct {
	Delegator columns.UserId `json:"delegator"`
	Delegee   columns.UserId `json:"delegee"`
	Objective models.Ovid    `json:"objective"`
}

type DelegateObjectiveResponse struct {
	Delid columns.DelegationId `json:"delid"`
}

// TODO: only last active delegee or the owner (if there is no delegee) can delegeate
// TODO: ensure there is no active collaboration
// POST
func (p *Public) DelegateObjective(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	_, err = p.sessions.WhoIs(&endpoints.WhoIsRequest{SessionToken: st})
	if err != nil {
		p.l.Println(fmt.Errorf("who is: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &DelegateObjectiveRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	panic("to implement") // TODO:

	bs := DelegateObjectiveResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
