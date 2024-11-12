package endpoints

import (
	"fmt"
	"logbook/internal/cookies"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type InviteCollaboratorsRequest struct {
	Collaborators []columns.UserId `json:"collaborators"`
}

type InviteCollaboratorsResponse struct {
	// TODO:
}

// TODO: check the inviter is owner; or the last delegee if there is any active delegation.
// TODO: check if the owner actually have right (connection) to send invites to invitees
// POST
func (p *Public) InviteCollaborators(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &InviteCollaboratorsRequest{}

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

	bs := InviteCollaboratorsResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
