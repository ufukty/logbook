package public

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type DelegateObjectiveRequest struct {
	Delegator columns.UserId
	Delegee   columns.UserId
	Objective models.Ovid
}

func (bq DelegateObjectiveRequest) validate() error {
	return validate.RequestFields(bq) // TODO: customize?
}

type DelegateObjectiveResponse struct {
	Delid columns.DelegationId
}

// TODO: only last active delegee or the owner (if there is no delegee) can delegeate
// TODO: ensure there is no active collaboration
func (e *Endpoints) DelegateObjective(w http.ResponseWriter, r *http.Request) {
	bq := &DelegateObjectiveRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		e.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	panic("to implement") // TODO:

	bs := DelegateObjectiveResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
