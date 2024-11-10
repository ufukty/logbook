package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type TagCreationRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Title        columns.TagTitle                      `json:"title"`
}

type TagCreationResponse struct {
	// TODO:
}

func (e *Endpoints) TagCreation(w http.ResponseWriter, r *http.Request) {
	bq := &TagCreationRequest{}

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

	bs := TagCreationResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
