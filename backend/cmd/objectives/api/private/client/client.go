package objectives

import (
	"logbook/cmd/objectives/api/private/endpoints"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
	"net/http"
)

type Client struct {
	lb         *balancer.LoadBalancer
	servicecfg api.Objectives
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config) *Client {
	return &Client{
		servicecfg: apicfg.Objectives,
		lb:         lb,
	}
}

func (c *Client) RockCreate(bq *endpoints.RockCreateRequest) (*http.Response, error) {
	return requests.BalancedSendRaw(c.lb, "", c.servicecfg.Private.RockCreate, bq)
}
