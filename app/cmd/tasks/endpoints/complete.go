package endpoints

import (
	"fmt"
	"log"
	"logbook/internal/web/reqs"
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

func (e *Endpoints) MarkComplete(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[MarkCompleteRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		log.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	panic("to implement") // TODO:

	bs := MarkCompleteResponse{} // TODO:
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
