package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/discovery/app"
	"logbook/internal/web/reqs"
	"logbook/models"
	"net/http"
	"net/url"
)

type RegisterInstanceRequest struct {
	Service models.Service `json:"service"`
	Address string         `json:"address"`
	Port    string         `json:"port"`
}

func (bq RegisterInstanceRequest) validate() error {
	u := fmt.Sprintf("%s:%s", bq.Address, bq.Port)
	_, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("declared address and port %q is invalid", u)
	}
	return nil
}

type RegisterInstanceResponse struct {
	InstanceId app.InstanceId `json:"instance-id"`
}

func (e *Endpoints) RegisterInstance(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[RegisterInstanceRequest](r)
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

	iid, err := e.a.RegisterInstance(bq.Service, app.Instance{
		Address: bq.Address,
		Port:    bq.Port,
	})
	if err != nil {
		log.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	bs := RegisterInstanceResponse{
		InstanceId: iid,
	}
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}

// func (bq *RegisterInstanceRequest) Send() (*RegisterInstanceResponse, error) {
// 	return reqs.Send[RegisterInstanceRequest, RegisterInstanceResponse](config., bq)
// }
