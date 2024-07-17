package models

import (
	"fmt"
	"net/url"
)

type Instance struct {
	Tls     bool   `json:"tls"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (i Instance) Url() url.URL {
	schema := "http"
	if i.Tls {
		schema = "https"
	}
	return url.URL{
		Scheme: schema,
		Host:   fmt.Sprintf(i.Address, i.Port),
	}
}
