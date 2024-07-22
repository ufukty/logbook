package servicereg

import (
	"fmt"
	"logbook/cmd/registry/endpoints"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
	"net/http"
	"path/filepath"
)

type Client struct {
	lb          *balancer.LoadBalancer
	servicepath string
	servicecfg  api.Discovery
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config, throughgateway bool) *Client {
	servicepath := apicfg.Internal.Services.Discovery.Path
	if throughgateway {
		servicepath = filepath.Join(apicfg.Internal.Path, servicepath)
	}
	return &Client{
		servicepath: servicepath,
		servicecfg:  apicfg.Internal.Services.Discovery,
		lb:          lb,
	}
}

func (c *Client) RegisterInstance(bq *endpoints.RegisterInstanceRequest) (*endpoints.RegisterInstanceResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.RegisterInstance.Path)
	bs := &endpoints.RegisterInstanceResponse{}
	err = requests.Send(url, c.servicecfg.Endpoints.RegisterInstance.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("requests.Send: %w", err)
	}
	return bs, err
}

func (c *Client) RecheckInstance(bq *endpoints.RecheckInstanceRequest) (*http.Response, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.RecheckInstance.Path)
	rs, err := requests.SendRaw(url, c.servicecfg.Endpoints.RecheckInstance.Method, bq)
	if err != nil {
		return nil, fmt.Errorf("requests.SendRaw: %w", err)
	}
	return rs, nil
}

func (c *Client) ListInstances(bq *endpoints.ListInstancesRequest) (*endpoints.ListInstancesResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.ListInstances.Path)
	bs := &endpoints.ListInstancesResponse{}
	err = requests.Send(url, c.servicecfg.Endpoints.ListInstances.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("requests.Send: %w", err)
	}
	return bs, err
}
