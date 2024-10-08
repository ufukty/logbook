package objectives

import (
	"logbook/config/api"
	"logbook/internal/web/balancer"
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
