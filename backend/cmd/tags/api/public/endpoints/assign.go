package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration/middlewares"
	"net/http"
)

type TagAssignRequest struct {
	// TODO:
}

func (bq TagAssignRequest) validate() error {
	panic("to implement") // TODO:
	return nil
}

type TagAssignResponse struct {
	// TODO:
}

func (e *Endpoints) TagAssign(id receptionist.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &TagAssignRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	panic("to implement") // TODO:

	bs := TagAssignResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}
