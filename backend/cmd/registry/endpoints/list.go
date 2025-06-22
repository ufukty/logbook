package endpoints

import (
	"fmt"
	"logbook/internal/web/serialize"
	"logbook/models"
	"net/http"
)

type ListInstancesRequest struct {
	Service models.Service `route:"service"`
}

type ListInstancesResponse struct {
	Instances []models.Instance `json:"instances"`
}

// GET
func (e *Endpoints) ListInstances(w http.ResponseWriter, r *http.Request) {
	bq := &ListInstancesRequest{}

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

	instances, err := e.a.ListInstances(bq.Service)
	if err != nil {
		e.l.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
	bs := ListInstancesResponse{instances}

	if err := bs.Write(w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}
