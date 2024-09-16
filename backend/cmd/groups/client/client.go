package objectives

import (
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
)

type Client struct {
	lb          *balancer.LoadBalancer
	servicepath string
	servicecfg  api.Objectives
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config) *Client {
	return &Client{
		servicepath: apicfg.Public.Services.Objectives.Path,
		servicecfg:  apicfg.Public.Services.Objectives,
		lb:          lb,
	}
}

func (c *Client) MarkComplete(bq *endpoints.MarkCompleteRequest) (*endpoints.MarkCompleteResponse, error) {
	return requests.BalancedSend(c.lb, c.servicepath, c.servicecfg.Endpoints.Mark, bq, &endpoints.MarkCompleteResponse{})
}

func (c *Client) CreateObjective(bq *endpoints.CreateObjectiveRequest) (*endpoints.CreateObjectiveResponse, error) {
	return requests.BalancedSend(c.lb, c.servicepath, c.servicecfg.Endpoints.Create, bq, &endpoints.CreateObjectiveResponse{})
}

func (c *Client) ReattachObjective(bq *endpoints.ReattachObjectiveRequest) (*endpoints.ReattachObjectiveResponse, error) {
	return requests.BalancedSend(c.lb, c.servicepath, c.servicecfg.Endpoints.Attach, bq, &endpoints.ReattachObjectiveResponse{})
}

func (c *Client) GetPlacementArray(bq *endpoints.GetPlacementArrayRequest) (*endpoints.GetPlacementArrayResponse, error) {
	return requests.BalancedSend(c.lb, c.servicepath, c.servicecfg.Endpoints.Placement, bq, &endpoints.GetPlacementArrayResponse{})
}
