package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/reception"
	"logbook/models"
	"net/http"
)

type ListInstancesRequest struct {
	Service models.Service `url:"service"`
}

type ListInstancesResponse struct {
	Instances []models.Instance `json:"instances"`
}

func (e *Endpoints) ListInstances(id reception.RequestId, store *reception.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &ListInstancesRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	instances, err := e.a.ListInstances(bq.Service)
	if err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("performing request: %w", err)
	}
	bs := ListInstancesResponse{instances}

	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}

// func (bq *ListInstancesRequest) Send() (*ListInstancesResponse, error) {
// 	return reqs.Send[ListInstancesRequest, ListInstancesResponse](config., bq)
// }
