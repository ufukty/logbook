package endpoints

import (
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type CreateGroupRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Name         columns.GroupName                     `json:"name"`
}

type CreateGroupResponse struct {
	GroupId columns.GroupId `json:"gid"`
}

// POST
func (e *Public) CreateGroup(w http.ResponseWriter, r *http.Request) {
	bq := &CreateGroupRequest{}

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

	panic("decide user id")

	gid, err := e.a.CreateGroup(r.Context(), app.CreateGroupParams{
		Actor:     "", // FIXME:
		GroupName: bq.Name,
	})
	if err != nil {
		e.l.Println(fmt.Errorf("CreateGroup: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateGroupResponse{
		GroupId: gid,
	}
	if err := bs.Write(w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		return
	}
}
