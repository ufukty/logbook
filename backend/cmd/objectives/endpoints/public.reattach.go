package endpoints

import (
	"fmt"
	"logbook/internal/cookies"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type ReattachObjectiveRequest struct {
	Subject   models.Ovid         `json:"subject"`
	NewParent columns.ObjectiveId `json:"new-parent"`
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

	panic("to implement") // TODO:

	bs := ReattachObjectiveResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		return
	}
}
