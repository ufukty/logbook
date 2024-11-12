package endpoints

import (
	"fmt"
	"logbook/cmd/account/endpoints"
	"logbook/cmd/groups/app"
	"logbook/internal/cookies"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type InviteMembersRequest struct {
	Gid              columns.GroupId   `route:"gid"`
	GroupTypeMembers []columns.GroupId `json:"user-type-members"`
	UserTypeMembers  []columns.UserId  `json:"group-type-members"`
}

// TODO: check the inviter is owner; or the last delegee if there is any active delegation.
// TODO: check if the owner actually have right (connection) to send invites to invitees
// POST
func (p *Public) InviteMembers(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	rp, err := p.accounts.WhoIs(&endpoints.WhoIsRequest{SessionToken: st})
	if err != nil {
		p.l.Println(fmt.Errorf("who is: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &InviteMembersRequest{}

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

	err = p.a.InviteMembers(r.Context(), app.InviteMembersParams{
		Actor:            rp.Uid,
		Gid:              bq.Gid,
		GroupTypeMembers: bq.GroupTypeMembers,
		UserTypeMembers:  bq.UserTypeMembers,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("a.InviteMembers: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
