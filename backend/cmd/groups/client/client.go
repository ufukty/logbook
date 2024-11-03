package objectives

import (
	"logbook/cmd/groups/endpoints"
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

func (c *Client) MembershipCheck(bq *endpoints.CheckMembershipRequest) (*endpoints.CheckMembershipResponse, error) {
	return requests.BalancedSend(c.lb, "", c.servicecfg.Private.GroupMembersCheck, bq, &endpoints.CheckMembershipResponse{})
}

func (c *Client) MembershipCheckEventual(bq *endpoints.CheckMembershipEventualRequest) (*endpoints.CheckMembershipEventualResponse, error) {
	return requests.BalancedSend(c.lb, "", c.servicecfg.Private.GroupMembersCheckEventual, bq, &endpoints.CheckMembershipEventualResponse{})
}
