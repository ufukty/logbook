package account

import (
	"logbook/config/api"
)

type Client struct {
	path   string
	config api.Account
}

func NewClient(apicfg api.Config) *Client {
	return &Client{
		path:   api.PathFromInternet(apicfg.Public.Services.Account),
		config: apicfg.Public.Services.Account,
	}
}

// func (c Client) Register(bq *endpoints.CreateUserRequest) (*endpoints.CreateUserResponse, error) {
// 	return reqs.SendTo[endpoints.CreateUserRequest, endpoints.CreateUserResponse](c.path, c.config.Endpoints.Create, bq)
// }
