package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type InviteCollaboratorsRequest struct {
	Collaborators []columns.UserId
}

func (bq InviteCollaboratorsRequest) validate() error {
	return validate.RequestFields(bq) // TODO: customize?
}

type InviteCollaboratorsResponse struct {
	// TODO:
}

// TODO: check the inviter is owner; or the last delegee if there is any active delegation.
// TODO: check if the owner actually have right (connection) to send invites to invitees
func (e *Endpoints) InviteCollaborators(w http.ResponseWriter, r *http.Request) {
	bq := &InviteCollaboratorsRequest{}

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

	bs := InviteCollaboratorsResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
