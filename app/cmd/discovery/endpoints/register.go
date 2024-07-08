package endpoints

import (
	"fmt"
	"log"
	"logbook/internal/web/reqs"
	"net/http"
)

type RegisterInstanceRequest struct {
	// TODO:
}

func (bq RegisterInstanceRequest) validate() error {
	panic("to implement") // TODO:
	return nil
}

type RegisterInstanceResponse struct {
	// TODO:
}

func (e *Endpoints) RegisterInstance(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[RegisterInstanceRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		log.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	panic("to implement") // TODO:

	bs := RegisterInstanceResponse{} // TODO:
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}

// func (bq *RegisterInstanceRequest) Send() (*RegisterInstanceResponse, error) {
// 	return reqs.Send[RegisterInstanceRequest, RegisterInstanceResponse](config., bq)
// }
