package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type InviteCollaboratorsRequest struct {
	SessionToken  requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Collaborators []columns.UserId                      `json:"collaborators"`
}

type InviteCollaboratorsResponse struct {
	// TODO:
}

// TODO: check the inviter is owner; or the last delegee if there is any active delegation.
// TODO: check if the owner actually have right (connection) to send invites to invitees
// POST
func (e *Public) InviteCollaborators(w http.ResponseWriter, r *http.Request) {
	bq := &InviteCollaboratorsRequest{}

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

	bs := InviteCollaboratorsResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
