package endpoints

import (
	"fmt"
	"log"
	"logbook/internal/web/requests"
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

func (e *Endpoints) TagCreation(w http.ResponseWriter, r *http.Request) {
	bq, err := requests.ParseRequest[TagCreationRequest](r)
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

	bs := TagCreationResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
