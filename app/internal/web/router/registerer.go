package router

import (
	"logbook/internal/web/api"
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
		return eps[i].Path > eps[j].Path
	})
	sort.Slice(eps, func(i, j int) bool {
		return strings.HasPrefix(string(eps[i].Path), string(eps[j].Path))
	})
	return eps
}
