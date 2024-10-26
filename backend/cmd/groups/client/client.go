package objectives

import (
	"logbook/cmd/groups/endpoints/private"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
)

type Client struct {
	lb         *balancer.LoadBalancer
	servicecfg api.Groups
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config) *Client {
	return &Client{
		servicecfg: apicfg.Groups,
		lb:         lb,
	}
}

func (c *Client) MembershipCheck(bq *private.MembershipCheckRequest) (*private.MembershipCheckResponse, error) {
	return requests.BalancedSend(c.lb, "", c.servicecfg.Private.MembershipCheck, bq, &private.MembershipCheckResponse{})
}
