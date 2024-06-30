package router

import (
	"logbook/config/api"
	"net/http"
	"sort"
	"strings"
)

type Endpoint struct {
	path        string
	method      string
	handler     http.Handler
	middlewares []http.Handler
}

type RequestId string

type RequestDetails struct {
	rid RequestId
}

type Handler func(d RequestDetails, w http.ResponseWriter, r *http.Request)

func sortEndpoints(eps []api.Endpoint) []api.Endpoint {
	sort.Slice(eps, func(i, j int) bool {
		return eps[i].GetPath() > eps[j].GetPath()
	})
	sort.Slice(eps, func(i, j int) bool {
		return strings.HasPrefix(string(eps[i].GetPath()), string(eps[j].GetPath()))
	})
	return eps
}
