package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/models"
	"net/http"
)

type ListInstancesRequest struct {
	Service models.Service `route:"service"`
}

type ListInstancesResponse struct {
	Instances []models.Instance `json:"instances"`
}

func (e *Endpoints) ListInstances(w http.ResponseWriter, r *http.Request) {
	bq := &ListInstancesRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	instances, err := e.a.ListInstances(bq.Service)
	if err != nil {
		e.l.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
	bs := ListInstancesResponse{instances}

	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}
