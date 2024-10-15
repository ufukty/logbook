package endpoints

import (
	"fmt"
	"logbook/cmd/registry/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/reception"
	"logbook/models"
	"net/http"
)

type RecheckInstanceRequest struct {
	Service    models.Service `json:"service"`
	InstanceId app.InstanceId `json:"instance-id"`
}

func (e *Endpoints) RecheckInstance(id reception.RequestId, store *reception.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &RecheckInstanceRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := e.a.RecheckInstance(bq.Service, bq.InstanceId); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("performing request: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// func (bq *RecheckInstanceRequest) Send() (*RecheckInstanceResponse, error) {
// 	return reqs.Send[RecheckInstanceRequest, RecheckInstanceResponse](config., bq)
// }
