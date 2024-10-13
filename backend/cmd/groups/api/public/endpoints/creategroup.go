package endpoints

import (
	"fmt"
	"logbook/cmd/groups/api/public/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration/middlewares"
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

func (e *Endpoints) CreateGroup(id receptionist.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &CreateGroupRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	panic("decide user id")

	gid, err := e.a.CreateGroup(r.Context(), app.CreateGroupParams{
		Actor:     "", // FIXME:
		GroupName: bq.Name,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("CreateGroup: %w", err)
	}

	bs := CreateGroupResponse{
		GroupId: gid,
	}
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}
