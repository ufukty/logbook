package objectives

import (
	"logbook/config/api"
	"logbook/internal/web/balancer"
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
