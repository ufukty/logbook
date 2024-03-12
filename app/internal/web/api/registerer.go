package api

import (
	"logbook/internal/web/logger"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
)

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
	l := logger.NewLogger("RouteRegisterer")
	return func(r *mux.Router) {
		r = r.UseEncodedPath()
		l.Println("Registering routes in order:")
		for _, ep := range sortEndpoints(maps.Keys(handlers)) {
			handler := handlers[ep]
			l.Printf("%s %s -> %p\n", ep.Method, ep.Path, handler)
			r.HandleFunc(string(ep.Path), handler).Methods(ep.Method)
		}
	}
}
