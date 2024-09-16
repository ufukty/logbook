package requests

import (
	"fmt"
	"logbook/config/api"
	"logbook/internal/utilities/urls"
	"logbook/internal/web/balancer"
	"net/http"
)

func BalancedSendRaw[Request any](lb *balancer.LoadBalancer, servicepath string, ep api.Endpoint, bq *Request) (*http.Response, error) {
	instance, err := lb.Next()
	if err != nil {
		return nil, fmt.Errorf("lb.Next: %w", err)
	}
	url := urls.Join(instance.String(), servicepath, ep.GetPath())
	rs, err := SendRaw(url, ep.GetMethod(), bq)
	if err != nil {
		return nil, fmt.Errorf("SendRaw: %w", err)
	}
	return rs, nil
}

func BalancedSend[Request, Response any](lb *balancer.LoadBalancer, servicepath string, ep api.Endpoint, bq *Request, bs *Response) (*Response, error) {
	instance, err := lb.Next()
	if err != nil {
		return nil, fmt.Errorf("lb.Next: %w", err)
	}
	url := urls.Join(instance.String(), servicepath, ep.GetPath())
	err = Send(url, ep.GetMethod(), bq, bs)
	if err != nil {
		return nil, fmt.Errorf("Send: %w", err)
	}
	return bs, nil
}
