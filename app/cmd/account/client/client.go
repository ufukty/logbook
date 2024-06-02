package account

import (
	"logbook/config/api"
)

type Client struct {
	path   api.Path
	config api.Account
}

func NewClient(apicfg api.Config) *Client {
	return &Client{
		path:   apicfg.Gateway.Path.Join(apicfg.Account.Path),
		config: apicfg.Account,
	}
}

// func (c Client) Register(bq *endpoints.CreateUserRequest) (*endpoints.CreateUserResponse, error) {
// 	return reqs.SendTo[endpoints.CreateUserRequest, endpoints.CreateUserResponse](c.path, c.config.Endpoints.Register, bq)
// }
