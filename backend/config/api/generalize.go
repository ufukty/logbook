package api

import "strings"

type Endpoint string

func (ep Endpoint) GetMethod() string {
	return strings.Split(string(ep), " ")[0]
}

func (ep Endpoint) GetPath() string {
	return strings.Split(string(ep), " ")[1]
}
