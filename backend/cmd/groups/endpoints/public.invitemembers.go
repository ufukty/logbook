package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"net/http"
)

type InviteMembersRequest struct {
	Gid              columns.GroupId  `route:"gid"`
	GroupTypeMembers columns.GroupIds `json:"user-type-members"`
	UserTypeMembers  columns.UserIds  `json:"group-type-members"`
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

	rp, err := p.sessions.WhoIs(&endpoints.WhoIsRequest{SessionToken: st})
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

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
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
