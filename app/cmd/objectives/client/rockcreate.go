package objectives

import (
	"fmt"
	"logbook/cmd/objectives/endpoints"
	"logbook/internal/utilities/urls"
	"logbook/internal/web/requests"
	"net/http"
)

func (c *Client) RockCreate(bq *endpoints.RockCreateRequest) (*http.Response, error) {
	instance, err := c.lb.Next()
	if err != nil {
		return nil, fmt.Errorf("LoadBalancer.Next: %w", err)
	}
	url := urls.Join(instance.String(), c.servicepath, c.servicecfg.Endpoints.RockCreate.Path)
	rs, err := requests.SendRaw(url, c.servicecfg.Endpoints.RockCreate.Method, bq)
	if err != nil {
		return nil, fmt.Errorf("requests.SendRaw: %w", err)
	}
	return rs, nil
}
