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

func (i Instance) Schema() string {
	if i.Tls {
		return "https"
	} else {
		return "http"
	}
}

func (i Instance) Url() url.URL {
	return url.URL{
		Scheme: i.Schema(),
		Host:   fmt.Sprintf(i.Address, i.Port),
	}
}

func (i Instance) String() string {
	return fmt.Sprintf("%s://%s:%d", i.Schema(), i.Address, i.Port)
}
