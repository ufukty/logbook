package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/groups/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type CreateGroupRequest struct {
	Name columns.GroupName `json:"name"`
}

func (bq CreateGroupRequest) validate() error {
	return validate.RequestFields(bq)
}

type CreateGroupResponse struct {
	GroupId columns.GroupId `json:"gid"`
}

func (e *Endpoints) CreateGroup(w http.ResponseWriter, r *http.Request) {
	bq := &CreateGroupRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		log.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	panic("decide user id")

	gid, err := e.app.CreateGroup(r.Context(), app.CreateGroupParams{
		Actor:     "", // FIXME:
		GroupName: bq.Name,
	})
	if err != nil {
		log.Println(fmt.Errorf("CreateGroup: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateGroupResponse{
		GroupId: gid,
	}
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
