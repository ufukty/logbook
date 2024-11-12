package endpoints

import (
	"fmt"
	"logbook/cmd/account/endpoints"
	"logbook/cmd/groups/app"
	"logbook/internal/cookies"
	"logbook/internal/web/validate"
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

	rp, err := p.accounts.WhoIs(&endpoints.WhoIsRequest{SessionToken: st})
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

	if err := validate.RequestFields(bq); err != nil {
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = p.a.RespondToUserInvite(r.Context(), app.RespondToUserInviteParams{
		Actor: rp.Uid,
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
