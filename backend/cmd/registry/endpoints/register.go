package endpoints

import (
	"fmt"
	"logbook/cmd/registry/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/reception"
	"logbook/models"
	"net/http"
	"net/url"
)

type RegisterInstanceRequest struct {
	Service models.Service `json:"service"`
	TLS     bool           `json:"tls"`
	Address string         `json:"address"`
	Port    int            `json:"port"`
}

func (bq RegisterInstanceRequest) validate() error {
	proto := "http"
	if bq.TLS {
		proto = "https"
	}
	u := fmt.Sprintf("%s://%s:%d", proto, bq.Address, bq.Port)
	_, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("declared address and port %q is invalid", u)
	}
	return nil
}

type RegisterInstanceResponse struct {
	InstanceId app.InstanceId `json:"instance-id"`
}

func (e *Endpoints) RegisterInstance(id reception.RequestId, store *reception.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &RegisterInstanceRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	iid, err := e.a.RegisterInstance(bq.Service, models.Instance{
		Tls:     bq.TLS,
		Address: bq.Address,
		Port:    bq.Port,
	})
	if err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("performing request: %w", err)
	}

	bs := RegisterInstanceResponse{
		InstanceId: iid,
	}
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("writing json response: %w", err)
	}

	return nil
}

// func (bq *RegisterInstanceRequest) Send() (*RegisterInstanceResponse, error) {
// 	return reqs.Send[RegisterInstanceRequest, RegisterInstanceResponse](config., bq)
// }
