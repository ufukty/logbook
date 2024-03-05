package account

import (
	"logbook/cmd/user/endpoints"
	"logbook/internal/web/api"
	"logbook/internal/web/reqs"
	"path/filepath"
)

type Client struct {
	path   string
	config api.Account
}

func NewClient(config api.Config) *Client {
	return &Client{
		path:   filepath.Join(string(config.Gateways.Public.Path), string(config.Gateways.Public.Services.Objectives.Path)),
		config: config.Gateways.Public.Services.Account,
	}
}

func (c Client) Register(bq *endpoints.CreateUserRequest) (*endpoints.CreateUserResponse, error) {
	return reqs.SendTo[endpoints.CreateUserRequest, endpoints.CreateUserResponse](c.path, c.config.Endpoints.Register, bq)
}
