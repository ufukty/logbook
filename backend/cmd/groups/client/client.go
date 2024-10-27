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

func (c *Client) MembershipCheck(bq *private.CheckMembershipRequest) (*private.CheckMembershipResponse, error) {
	return requests.BalancedSend(c.lb, "", c.servicecfg.Private.GroupMembersCheck, bq, &private.CheckMembershipResponse{})
}

func (c *Client) MembershipCheckEventual(bq *private.CheckMembershipEventualRequest) (*private.CheckMembershipEventualResponse, error) {
	return requests.BalancedSend(c.lb, "", c.servicecfg.Private.GroupMembersCheckEventual, bq, &private.CheckMembershipEventualResponse{})
}
