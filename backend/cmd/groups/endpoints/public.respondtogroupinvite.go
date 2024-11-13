package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type RespondToGroupInviteRequest struct {
	Gid        columns.GroupId           `json:"ginvid"`
	Ginvid     columns.GroupInviteId     `json:"ginvid"`
	Response   transports.InviteResponse `json:"response"`
	MemberType transports.MemberType     `json:"member-type"`
}

// POST
func (p *Public) RespondToGroupInvite(w http.ResponseWriter, r *http.Request) {
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

	bq := &RespondToGroupInviteRequest{}

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

	err = p.a.RespondToGroupInvite(r.Context(), app.RespondToGroupInviteParams{
		Actor:      rp.Uid,
		Behalf:     bq.Gid,
		Ginvid:     bq.Ginvid,
		MemberType: bq.MemberType,
		Response:   bq.Response,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("app.RespondToGroupInvite: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
