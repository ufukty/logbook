package account

import (
	"logbook/cmd/account/api/private/endpoints"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
)

type Client struct {
	lb         *balancer.LoadBalancer
	servicecfg api.Account
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config) *Client {
	return &Client{
		servicecfg: apicfg.Account,
		lb:         lb,
	}
}

func (c *Client) WhoIs(bq *endpoints.WhoIsRequest) (*endpoints.WhoIsResponse, error) {
	return requests.BalancedSend(c.lb, "", c.servicecfg.Private.WhoIs, bq, &endpoints.WhoIsResponse{})
}
