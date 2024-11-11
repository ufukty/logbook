package endpoints

import (
	"fmt"
	"logbook/cmd/registry/app"
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
	if err := bq.Service.Validate(); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	return nil
}

type RegisterInstanceResponse struct {
	InstanceId app.InstanceId `json:"instance-id"`
}

// POST
func (e *Endpoints) RegisterInstance(w http.ResponseWriter, r *http.Request) {
	bq := &RegisterInstanceRequest{}

	if err := bq.Parse(r); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		e.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	iid, err := e.a.RegisterInstance(bq.Service, models.Instance{
		Tls:     bq.TLS,
		Address: bq.Address,
		Port:    bq.Port,
	})
	if err != nil {
		e.l.Println(fmt.Errorf("performing request: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	bs := RegisterInstanceResponse{
		InstanceId: iid,
	}
	if err := bs.Write(w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}
