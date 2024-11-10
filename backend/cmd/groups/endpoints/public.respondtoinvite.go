package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"logbook/models/incoming"
	"net/http"
)

type RespondToInviteRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Ginvid       columns.GroupInviteId                 `json:"ginvid"`
	Response     incoming.InviteResponse               `json:"response"`
	MemberType   incoming.MemberType                   `json:"member-type"`
}

type RespondToInviteResponse struct {
	// TODO:
}

// POST
func (e *Public) RespondToInvite(w http.ResponseWriter, r *http.Request) {
	bq := &RespondToInviteRequest{}

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

	bs := RespondToInviteResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
