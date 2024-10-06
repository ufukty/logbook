package client

import (
	"fmt"
	"logbook/cmd/tags/endpoints"
	"logbook/config/api"
	"logbook/internal/utils/urls"
	"logbook/internal/web/balancer"
	"logbook/internal/web/requests"
)

type Client struct {
	lb          *balancer.LoadBalancer
	servicepath string
	servicecfg  api.Tags
}

func NewClient(lb *balancer.LoadBalancer, apicfg *api.Config) *Client {
	return &Client{
		servicepath: apicfg.Public.Services.Tags.Path,
		servicecfg:  apicfg.Public.Services.Tags,
		lb:          lb,
	}
}

func (c *Client) TagAssign(bq *endpoints.TagAssignRequest) (*endpoints.TagAssignResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := urls.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.Assign.Path)
	bs := &endpoints.TagAssignResponse{}
	err = requests.Send(url, c.servicecfg.Endpoints.Assign.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("requests.Send: %w", err)
	}
	return bs, err
}

func (c *Client) TagCreation(bq *endpoints.TagCreationRequest) (*endpoints.TagCreationResponse, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := urls.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.Creation.Path)
	bs := &endpoints.TagCreationResponse{}
	err = requests.Send(url, c.servicecfg.Endpoints.Creation.Method, bq, bs)
	if err != nil {
		return nil, fmt.Errorf("requests.Send: %w", err)
	}
	return bs, err
}
