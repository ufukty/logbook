package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/api/public/middlewares"
	"logbook/cmd/objectives/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type CreateObjectiveRequest struct {
	Parent  models.Ovid      `json:"parent"`
	Content ObjectiveContent `json:"content"`
}

func (ct CreateObjectiveRequest) validate() error {
	return validate.RequestFields(ct)
}

type CreateObjectiveResponse struct {
	Oid columns.ObjectiveId `json:"oid"`
}

// TODO: Check user input for script tags in order to prevent XSS attempts
func (e *Endpoints) CreateObjective(rid receptionist.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &CreateObjectiveRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	obj, err := e.a.CreateSubtask(r.Context(), app.CreateSubtaskParams{
		Parent:  bq.Parent,
		Content: string(bq.Content),
		Creator: "00000000-0000-0000-0000-000000000000", // TODO: check auth header
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("app.createObjective: %w", err)
	}

	bs := CreateObjectiveResponse{
		Oid: obj.Oid,
	}
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}
