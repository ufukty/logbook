package endpoints

import (
	"fmt"
	"logbook/cmd/registry/models/scalars"
	"logbook/internal/web/serialize"
	"logbook/models"
	"net/http"
)

type RecheckInstanceRequest struct {
	Service    models.Service     `json:"service"`
	InstanceId scalars.InstanceId `json:"instance-id"`
}

// POST
func (e *Endpoints) RecheckInstance(w http.ResponseWriter, r *http.Request) {
	bq := &RecheckInstanceRequest{}

	if err := bq.Parse(r); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			e.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	if err := e.a.RecheckInstance(bq.Service, bq.InstanceId); err != nil {
		e.l.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
