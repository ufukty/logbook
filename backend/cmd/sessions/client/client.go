// Code generated by gohandlers v0.27.3. DO NOT EDIT.

package sessions

import (
	"fmt"
	"logbook/cmd/sessions/endpoints"
	"net/http"
)

type Pool interface {
	Host() (string, error)
}

type Client struct {
	p Pool
}

func NewClient(p Pool) *Client {
	return &Client{p: p}
}

func (c *Client) Login(bq *endpoints.LoginRequest) (*http.Response, error) {
	h, err := c.p.Host()
	if err != nil {
		return nil, fmt.Errorf("Host: %w", err)
	}
	rq, err := bq.Build(h)
	if err != nil {
		return nil, fmt.Errorf("Build: %w", err)
	}
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("Do: %w", err)
	}
	if rs.StatusCode != http.StatusOK {
		return rs, fmt.Errorf("non-200 status code: %d (%s)", rs.StatusCode, http.StatusText(rs.StatusCode))
	}
	return rs, nil
}

func (c *Client) SaveCredentials(bq *endpoints.SaveCredentialsRequest) (*http.Response, error) {
	h, err := c.p.Host()
	if err != nil {
		return nil, fmt.Errorf("Host: %w", err)
	}
	rq, err := bq.Build(h)
	if err != nil {
		return nil, fmt.Errorf("Build: %w", err)
	}
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("Do: %w", err)
	}
	if rs.StatusCode != http.StatusOK {
		return rs, fmt.Errorf("non-200 status code: %d (%s)", rs.StatusCode, http.StatusText(rs.StatusCode))
	}
	return rs, nil
}

func (c *Client) WhoIs(bq *endpoints.WhoIsRequest) (*endpoints.WhoIsResponse, error) {
	h, err := c.p.Host()
	if err != nil {
		return nil, fmt.Errorf("Host: %w", err)
	}
	rq, err := bq.Build(h)
	if err != nil {
		return nil, fmt.Errorf("Build: %w", err)
	}
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("Do: %w", err)
	}
	if rs.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d (%s)", rs.StatusCode, http.StatusText(rs.StatusCode))
	}
	bs := &endpoints.WhoIsResponse{}
	err = bs.Parse(rs)
	if err != nil {
		return nil, fmt.Errorf("Parse: %w", err)
	}
	return bs, nil
}
