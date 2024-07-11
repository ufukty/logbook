package objectives

import (
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/requests"
)

type Client struct {
	servicepath string
	servicecfg  api.Objectives
}

func NewClient(apicfg api.Config) *Client {
	return &Client{
		servicepath: api.PathFromInternet(apicfg.Public.Services.Objectives),
		servicecfg:  apicfg.Public.Services.Objectives,
	}
}

func (c Client) CreateObjective(bq *endpoints.CreateTaskRequest) (*endpoints.CreateTaskResponse, error) {
	return requests.Send[endpoints.CreateTaskRequest, endpoints.CreateTaskResponse](
		api.PathFromInternet(c.servicecfg.Endpoints.Create), c.servicecfg.Endpoints.Create.Method, bq,
	)
}
