package objectives

import (
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/reqs"
	"path/filepath"
)

type Client struct {
	path   string
	config api.Objectives
}

func NewClient(config api.Config) *Client {
	return &Client{
		path:   filepath.Join(string(config.Gateways.Public.Path), string(config.Gateways.Public.Services.Objectives.Path)),
		config: config.Gateways.Public.Services.Objectives,
	}
}

func (c Client) CreateObjective(bq *endpoints.CreateTaskRequest) (*endpoints.CreateTaskResponse, error) {
	return reqs.SendTo[endpoints.CreateTaskRequest, endpoints.CreateTaskResponse](c.path, c.config.Endpoints.Create, bq)
}
