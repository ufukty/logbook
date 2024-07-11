package endpoints

import (
	"fmt"
	"log"
	"logbook/internal/web/requests"
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

func (e *Endpoints) ReattachObjective(w http.ResponseWriter, r *http.Request) {
	bq, err := requests.ParseRequest[ReattachObjectiveRequest](r)
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

	bs := ReattachObjectiveResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
