package registry

import (
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
	servicecfg  api.Registry
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config, throughgateway bool) *Client {
	servicepath := apicfg.Internal.Services.Registry.Path
	if throughgateway {
		servicepath = filepath.Join(apicfg.Internal.Path, servicepath)
	}
	return &Client{
		servicepath: servicepath,
		servicecfg:  apicfg.Internal.Services.Registry,
		lb:          lb,
	}
}

func (c *Client) RegisterInstance(bq *endpoints.RegisterInstanceRequest) (*endpoints.RegisterInstanceResponse, error) {
	return requests.BalancedSend(c.lb, c.servicepath, c.servicecfg.Endpoints.RegisterInstance, bq, &endpoints.RegisterInstanceResponse{})
}

func (c *Client) RecheckInstance(bq *endpoints.RecheckInstanceRequest) (*http.Response, error) {
	return requests.BalancedSendRaw(c.lb, c.servicepath, c.servicecfg.Endpoints.RecheckInstance, bq)
}

func (c *Client) ListInstances(bq *endpoints.ListInstancesRequest) (*endpoints.ListInstancesResponse, error) {
	return requests.BalancedSend(c.lb, c.servicepath, c.servicecfg.Endpoints.ListInstances, bq, &endpoints.ListInstancesResponse{})
}
