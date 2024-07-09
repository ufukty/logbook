package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/discovery/app"
	"logbook/internal/web/reqs"
	"net/http"
)

type RecheckInstanceRequest struct {
	InstanceId app.InstanceId `json:"instance-id"`
}

func (e *Endpoints) RecheckInstance(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[RecheckInstanceRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err = e.a.RecheckInstance(bq.InstanceId); err != nil {
		log.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// func (bq *RecheckInstanceRequest) Send() (*RecheckInstanceResponse, error) {
// 	return reqs.Send[RecheckInstanceRequest, RecheckInstanceResponse](config., bq)
// }
