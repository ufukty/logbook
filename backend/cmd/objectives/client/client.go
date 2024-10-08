package objectives

import (
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
	"net/http"
)

type Client struct {
	lb          *balancer.LoadBalancer
	servicepath string
	servicecfg  api.Objectives
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config) *Client {
	return &Client{
		servicepath: apicfg.Internal.Services.Objectives.Path,
		servicecfg:  apicfg.Internal.Services.Objectives,
		lb:          lb,
	}
}

func (c *Client) RockCreate(bq *endpoints.RockCreateRequest) (*http.Response, error) {
	return requests.BalancedSendRaw(c.lb, c.servicepath, c.servicecfg.Endpoints.RockCreate, bq)
}
