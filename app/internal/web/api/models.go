package api

import (
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
)

type Domain struct {
	Protocol string // eg. http, https
	Domain   string // eg. logbook.app
	Port     string // eg. 8080
}

type Gateway struct {
	// Root    Domain
	Listens string // eg. /api/v1.0.0
}

type Service struct {
	Gateway Gateway
	Listens string
}

// func (d Domain) Url() string {
// 	return fmt.Sprintf("%s://%s:%s", d.Protocol, d.Domain, d.Port)
// }

// func (g Gateway) Url() string {
// 	return filepath.Join(g.Root.Url(), g.Listens)
// }

// func (s Service) Url() string {
// 	return filepath.Join(s.Gateway.Url(), s.Listens)
// }

// // returns the path that is supposed to get trimmed at gateway from each request's Host header
// func (s Service) PathByGateway() string {
// 	return filepath.Join(s.Gateway.Listens, s.Listens)
// }

// func (e Endpoint) Url() string {
// 	return filepath.Join(e.Service.Url(), e.Listens)
// }

func (e Endpoint) Url(gateway, service Path) string {
	return filepath.Join(string(gateway), string(service), string(e.Path))
}

// func checkPrefix(a Endpoint, b Endpoint) bool {
// 	return strings.HasPrefix(a.Url(), b.Url())
// }

func sortEndpoints(eps []Endpoint) []Endpoint {
	sort.Slice(eps, func(i, j int) bool {
		return eps[i].Path > eps[j].Path
	})
	sort.Slice(eps, func(i, j int) bool {
		return strings.HasPrefix(string(eps[i].Path), string(eps[j].Path))
	})
	return eps
}

func RouteRegisterer(handlers map[Endpoint]http.HandlerFunc) func(*mux.Router) {
	return func(r *mux.Router) {
		r = r.UseEncodedPath()
		for _, ep := range sortEndpoints(maps.Keys(handlers)) {
			var handler = handlers[ep]
			log.Printf("Registering route: %-6s %s\n", ep.Method, ep.Path)
			r.HandleFunc(string(ep.Path), handler).Methods(ep.Method)
		}
	}
}
