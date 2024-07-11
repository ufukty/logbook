package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/discovery/app"
	"logbook/internal/web/requests"
	"logbook/models"
	"net/http"
)

type ListInstancesRequest struct {
	Service models.Service `url:"service"`
}

type ListInstancesResponse []app.Instance

func (e *Endpoints) ListInstances(w http.ResponseWriter, r *http.Request) {
	bq, err := requests.ParseRequest[ListInstancesRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	instances, err := e.a.ListInstances(bq.Service)
	bs := ListInstancesResponse(instances)

	if err := requests.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}

// func (bq *ListInstancesRequest) Send() (*ListInstancesResponse, error) {
// 	return reqs.Send[ListInstancesRequest, ListInstancesResponse](config., bq)
// }
