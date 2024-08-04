package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/objectives/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models"
	"net/http"
)

type CreateObjectiveRequest struct {
	Parent  models.Ovid      `json:"parent"`
	Content ObjectiveContent `json:"content"`
}

func (ct CreateObjectiveRequest) validate() error {
	return validate.RequestFields(ct)
}

type CreateObjectiveResponse struct{}

// TODO: Check user input for script tags in order to prevent XSS attempts
func (e *Endpoints) CreateObjective(w http.ResponseWriter, r *http.Request) {
	bq := &CreateObjectiveRequest{}

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

	err := e.app.CreateSubtask(r.Context(), app.CreateSubtaskParams{
		Parent:  bq.Parent,
		Content: string(bq.Content),
		Creator: "00000000-0000-0000-0000-000000000000", // TODO: check auth header
	})
	if err != nil {
		log.Println(fmt.Errorf("app.createObjective: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateObjectiveResponse{}
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
