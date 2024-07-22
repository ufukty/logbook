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

func (c *Client) CreateObjective(bq *endpoints.CreateTaskRequest) (*endpoints.CreateTaskResponse, error) {
	bs := &endpoints.CreateTaskResponse{}
	next, err := c.lb.Next()
	url := filepath.Join(next.String(), c.servicepath, c.servicecfg.Endpoints.Create.Path)
	err = requests.Send(url, c.servicecfg.Endpoints.Create.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}
	return bs, nil
}
