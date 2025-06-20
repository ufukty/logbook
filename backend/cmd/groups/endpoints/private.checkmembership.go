package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/models/columns"
	"net/http"
)

type CheckMembershipRequest struct {
	Uid columns.UserId  `route:"uid"`
	Gid columns.GroupId `route:"gid"`
}

type CheckMembershipResponse struct {
	Membership bool `json:"membership"`
}

// GET
func (e *Private) CheckMembership(w http.ResponseWriter, r *http.Request) {
	bq := &CheckMembershipRequest{}

	if err := bq.Parse(r); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		e.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	membership, err := e.a.CheckMembership(r.Context(), app.CheckMembershipParams{
		Uid:      bq.Uid,
		Gid:      bq.Gid,
		Eventual: false,
	})
	if err != nil {
		e.l.Println(fmt.Errorf("a.CheckMembership: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CheckMembershipResponse{
		Membership: membership,
	}
	if err := bs.Write(w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
