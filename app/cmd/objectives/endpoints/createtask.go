package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/objectives/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"net/http"
)

type CreateTaskRequest struct {
	Parent  app.Ovid         `json:"parent"`
	Content ObjectiveContent `json:"content"`
}

func (ct CreateTaskRequest) validate() error {
	return validate.RequestFields(ct)
}

type CreateTaskResponse struct {
	Update []app.Ovid `json:"update"`
}

// TODO: Check user input for script tags in order to prevent XSS attempts
func (e *Endpoints) CreateTask(w http.ResponseWriter, r *http.Request) {
	bq := &CreateTaskRequest{}

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

	o, err := e.app.CreateObjective(r.Context(), app.CreateObjectiveAction{
		Parent:  bq.Parent,
		Content: string(bq.Content),
		Creator: "00000000-0000-0000-0000-000000000000", // TODO: check auth header
	})
	if err != nil {
		log.Println(fmt.Errorf("app.createObjective: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateTaskResponse{
		Update: o,
	}
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
