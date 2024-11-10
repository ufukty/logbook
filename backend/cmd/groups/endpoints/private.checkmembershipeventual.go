package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type CheckMembershipEventualRequest struct {
	Uid columns.UserId  `route:"uid"`
	Gid columns.GroupId `route:"gid"`
}

type CheckMembershipEventualResponse struct {
	Membership bool `json:"membership"`
}

// GET
func (e *Private) CheckMembershipEventual(w http.ResponseWriter, r *http.Request) {
	bq := &CheckMembershipEventualRequest{}

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

	membership, err := e.a.CheckMembership(r.Context(), app.CheckMembershipParams{
		Uid:      bq.Uid,
		Gid:      bq.Gid,
		Eventual: true,
	})
	if err != nil {
		e.l.Println(fmt.Errorf("a.CheckMembership: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CheckMembershipEventualResponse{
		Membership: membership,
	}
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
