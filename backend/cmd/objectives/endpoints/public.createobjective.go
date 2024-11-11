package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type CreateObjectiveRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Parent       models.Ovid                           `json:"parent"`
	Content      columns.ObjectiveContent              `json:"content"`
}

type CreateObjectiveResponse struct {
	Oid columns.ObjectiveId `json:"oid"`
}

// TODO: Check user input for script tags in order to prevent XSS attempts
// POST
func (e *Public) CreateObjective(w http.ResponseWriter, r *http.Request) {
	bq := &CreateObjectiveRequest{}

	if err := bq.Parse(r); err != nil {
		fmt.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := validate.RequestFields(bq); err != nil {
		fmt.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	obj, err := e.a.CreateSubtask(r.Context(), app.CreateSubtaskParams{
		Parent:  bq.Parent,
		Content: bq.Content,
		Creator: "00000000-0000-0000-0000-000000000000", // TODO: check auth header
	})
	if err != nil {
		fmt.Println(fmt.Errorf("app.createObjective: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateObjectiveResponse{
		Oid: obj.Oid,
	}
	if err := bs.Write(w); err != nil {
		fmt.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
