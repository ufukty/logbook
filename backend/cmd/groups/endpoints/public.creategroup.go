package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/models/columns"
	"net/http"
)

type CreateGroupRequest struct {
	Name columns.GroupName `json:"name"`
}

type CreateGroupResponse struct {
	GroupId columns.GroupId `json:"gid"`
}

// POST
func (p *Public) CreateGroup(w http.ResponseWriter, r *http.Request) {
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

	bq := &CreateGroupRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	gid, err := p.a.CreateGroup(r.Context(), app.CreateGroupParams{
		Actor:     rp.Uid,
		GroupName: bq.Name,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("CreateGroup: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateGroupResponse{
		GroupId: gid,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		return
	}
}
