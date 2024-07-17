package discovery

import (
	"logbook/cmd/discovery/endpoints"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
)

type Client struct {
	lb *balancer.LoadBalancer
}

func NewClient(lb *balancer.LoadBalancer) *Client {
	return &Client{
		lb: lb,
	}
}

func (c Client) RegisterInstance(bq *endpoints.RegisterInstanceRequest) (*endpoints.RegisterInstanceResponse, error) {
	return requests.Send[endpoints.RegisterInstanceRequest, endpoints.RegisterInstanceResponse](c.lb.Next(), bq)
}

func (c Client) RecheckInstance(bq *endpoints.RecheckInstanceRequest) (*endpoints.RecheckInstanceResponse, error) {
	return requests.Send[endpoints.RecheckInstanceRequest, endpoints.RecheckInstanceResponse](c.lb.Next(), bq)
}

func (c Client) ListInstances(bq *endpoints.ListInstancesRequest) (*endpoints.ListInstancesResponse, error) {
	return requests.Send[endpoints.ListInstancesRequest, endpoints.ListInstancesResponse](c.lb.Next(), bq)
}
