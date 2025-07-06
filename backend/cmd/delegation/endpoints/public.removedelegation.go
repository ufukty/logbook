package endpoints

import (
	"fmt"
	"logbook/internal/cookies"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"net/http"
)

type RemoveDelegationRequest struct {
	Delid columns.DelegationId `json:"delid"`
}

type RemoveDelegationResponse struct {
	// TODO:
}

// POST
func (p *Public) RemoveDelegation(w http.ResponseWriter, r *http.Request) {
	_, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &RemoveDelegationRequest{}

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

	bs := RemoveDelegationResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
