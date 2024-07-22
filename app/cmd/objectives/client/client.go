package objectives

import (
	"fmt"
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
	"path/filepath"
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

func (c *Client) CreateObjective(bq *endpoints.CreateObjectiveRequest) (*endpoints.CreateObjectiveResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.Create.Path)
	bs := &endpoints.CreateObjectiveResponse{}
	err = requests.Send(url, c.servicecfg.Endpoints.Create.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("requests.Send: %w", err)
	}
	return bs, err
}

func (c *Client) ReattachObjective(bq *endpoints.ReattachObjectiveRequest) (*endpoints.ReattachObjectiveResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := filepath.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.Attach.Path)
	bs := &endpoints.ReattachObjectiveResponse{}
	err = requests.Send(url, c.servicecfg.Endpoints.Attach.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("requests.Send: %w", err)
	}
	return bs, err
}
