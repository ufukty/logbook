package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type DelegateObjectiveRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Delegator    columns.UserId                        `json:"delegator"`
	Delegee      columns.UserId                        `json:"delegee"`
	Objective    models.Ovid                           `json:"objective"`
}

type DelegateObjectiveResponse struct {
	Delid columns.DelegationId `json:"delid"`
}

// TODO: only last active delegee or the owner (if there is no delegee) can delegeate
// TODO: ensure there is no active collaboration
// POST
func (e *Public) DelegateObjective(w http.ResponseWriter, r *http.Request) {
	bq := &DelegateObjectiveRequest{}

	if err := bq.Parse(r); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := validate.RequestFields(bq); err != nil {
		e.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	panic("to implement") // TODO:

	bs := DelegateObjectiveResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
