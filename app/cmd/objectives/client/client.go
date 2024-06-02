package objectives

import (
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/reqs"
)

type Client struct {
	path   api.Path
	config api.Objectives
}

func NewClient(apicfg api.Config) *Client {
	return &Client{
		path:   apicfg.Gateway.Path.Join(apicfg.Objectives.Path),
		config: apicfg.Objectives,
	}
}

func (c Client) CreateObjective(bq *endpoints.CreateTaskRequest) (*endpoints.CreateTaskResponse, error) {
	return reqs.SendTo[endpoints.CreateTaskRequest, endpoints.CreateTaskResponse](string(c.path), c.config.Endpoints.Create, bq)
}
