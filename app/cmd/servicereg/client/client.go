package servicereg

import (
	"fmt"
	"logbook/cmd/servicereg/endpoints"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
	"net/http"
	"path/filepath"
)

type Client struct {
	lb          *balancer.LoadBalancer
	servicepath string
	servicecfg  api.Internal
}

// servicepath => https://ip:port/(servicepath)/endpoint, which changes when gateway used
func NewClient(apicfg *api.Config, lb *balancer.LoadBalancer, servicepath string) *Client {
	return &Client{
		servicepath: servicepath,
		servicecfg:  apicfg.Internal,
		lb:          lb,
	}
}

func (c *Client) RegisterInstance(bq *endpoints.RegisterInstanceRequest) (*endpoints.RegisterInstanceResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("getting the next service instance: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Services.Discovery.Endpoints.RecheckInstance.Path)
	bs := &endpoints.RegisterInstanceResponse{}
	err = requests.Send(url, c.servicecfg.Services.Discovery.Endpoints.RecheckInstance.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("sending over bindings: %w", err)
	}
	return bs, err
}

func (c *Client) RecheckInstance(bq *endpoints.RecheckInstanceRequest) (*http.Response, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("loadbalancer: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Services.Discovery.Endpoints.RecheckInstance.Path)
	return requests.SendRaw(url, c.servicecfg.Services.Discovery.Endpoints.RecheckInstance.Method, bq)
}

func (c *Client) ListInstances(bq *endpoints.ListInstancesRequest) (*endpoints.ListInstancesResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("getting the next service instance: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Services.Discovery.Endpoints.RecheckInstance.Path)
	bs := &endpoints.ListInstancesResponse{}
	err = requests.Send(url, c.servicecfg.Services.Discovery.Endpoints.RecheckInstance.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("sending over bindings: %w", err)
	}
	return bs, err
}
