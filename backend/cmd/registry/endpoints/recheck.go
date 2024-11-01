package endpoints

import (
	"fmt"
	"logbook/cmd/registry/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
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
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := validate.RequestFields(bq); err != nil {
		e.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := e.a.RecheckInstance(bq.Service, bq.InstanceId); err != nil {
		e.l.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
