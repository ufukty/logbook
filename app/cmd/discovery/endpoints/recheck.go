package endpoints

import (
	"fmt"
	"log"
	"logbook/internal/web/reqs"
	"net/http"
)

type RecheckInstanceRequest struct {
	// TODO:
}

func (bq RecheckInstanceRequest) validate() error {
	panic("to implement") // TODO:
	return nil
}

type RecheckInstanceResponse struct {
	// TODO:
}

func (e *Endpoints) RecheckInstance(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[RecheckInstanceRequest](r)
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

	bs := RecheckInstanceResponse{} // TODO:
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}

// func (bq *RecheckInstanceRequest) Send() (*RecheckInstanceResponse, error) {
// 	return reqs.Send[RecheckInstanceRequest, RecheckInstanceResponse](config., bq)
// }
