package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type GetPlacementArrayRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Root         models.Ovid                           `route:"root"`
}

type GetPlacementArrayResponse struct {
	Items []database.Objective `json:"items"`
}

// GET
func (e *Public) GetPlacementArray(w http.ResponseWriter, r *http.Request) {
	bq := &GetPlacementArrayRequest{}

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

	panic("to implement") // TODO:

	bs := GetPlacementArrayResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
