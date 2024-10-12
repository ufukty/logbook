package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/pipelines"
	"logbook/internal/web/router/pipelines/middlewares"
	"net/http"
)

type MarkCompleteRequest struct {
	// TODO:
}

func (bq MarkCompleteRequest) validate() error {
	panic("to implement") // TODO:
	return nil
}

type MarkCompleteResponse struct {
	// TODO:
}

func (e *Endpoints) MarkComplete(rid pipelines.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &MarkCompleteRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	panic("to implement") // TODO:

	bs := MarkCompleteResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}
