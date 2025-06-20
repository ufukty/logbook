package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/internal/web/serialize"
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
func (p *Private) CheckMembershipEventual(w http.ResponseWriter, r *http.Request) {
	bq := &CheckMembershipEventualRequest{}

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

	membership, err := p.a.CheckMembership(r.Context(), app.CheckMembershipParams{
		Uid:      bq.Uid,
		Gid:      bq.Gid,
		Eventual: true,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("a.CheckMembership: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CheckMembershipEventualResponse{
		Membership: membership,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
