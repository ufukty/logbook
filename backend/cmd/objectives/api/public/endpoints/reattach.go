package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/reception"
	"logbook/internal/web/router/reception/middlewares"
	"net/http"
)

type ReattachObjectiveRequest struct {
	// TODO:
}

func (bq ReattachObjectiveRequest) validate() error {
	panic("to implement") // TODO:
	return nil
}

type ReattachObjectiveResponse struct {
	// TODO:
}

func (e *Endpoints) ReattachObjective(rid reception.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &ReattachObjectiveRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	panic("to implement") // TODO:

	bs := ReattachObjectiveResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}
