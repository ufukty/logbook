package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/registry/app"
	"logbook/internal/web/requests"
	"logbook/models"
	"net/http"
)

type RecheckInstanceRequest struct {
	Service    models.Service `json:"service"`
	InstanceId app.InstanceId `json:"instance-id"`
}

func (e *Endpoints) RecheckInstance(w http.ResponseWriter, r *http.Request) {
	bq := &RecheckInstanceRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := e.a.RecheckInstance(bq.Service, bq.InstanceId); err != nil {
		log.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// func (bq *RecheckInstanceRequest) Send() (*RecheckInstanceResponse, error) {
// 	return reqs.Send[RecheckInstanceRequest, RecheckInstanceResponse](config., bq)
// }
