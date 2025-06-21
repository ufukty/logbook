package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type RespondToUserInviteRequest struct {
	Ginvid     columns.GroupInviteId     `json:"ginvid"`
	Response   transports.InviteResponse `json:"response"`
	MemberType transports.MemberType     `json:"member-type"`
}

// POST
func (p *Public) RespondToUserInvite(w http.ResponseWriter, r *http.Request) {
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

	bq := &RespondToUserInviteRequest{}

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
	err = p.a.RespondToUserInvite(r.Context(), app.RespondToUserInviteParams{
		Actor:      rp.Uid,
		Ginvid:     bq.Ginvid,
		MemberType: bq.MemberType,
		Response:   bq.Response,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("app.RespondToUserInvite: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
