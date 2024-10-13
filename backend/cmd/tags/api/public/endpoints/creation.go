package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration/middlewares"
	"net/http"
)

type TagCreationRequest struct {
	// TODO:
}

func (bq TagCreationRequest) validate() error {
	panic("to implement") // TODO:
	return nil
}

type TagCreationResponse struct {
	// TODO:
}

func (e *Endpoints) TagCreation(id receptionist.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &TagCreationRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	panic("to implement") // TODO:

	bs := TagCreationResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}
