package router

import (
	"fmt"
	"logbook/config/api"
	"logbook/internal/logger"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
)

func sortEndpoints(eps []api.Endpoint) []api.Endpoint {
	sort.Slice(eps, func(i, j int) bool {
		return eps[i].GetPath() > eps[j].GetPath()
	})
	sort.Slice(eps, func(i, j int) bool {
		return strings.HasPrefix(string(eps[i].GetPath()), string(eps[j].GetPath()))
	})
	return eps
}

type EndpointDetails struct {
	Handler http.HandlerFunc
	Cors    func(http.HandlerFunc) http.HandlerFunc // optional
}

func registerer(r *mux.Router, details map[api.Endpoint]EndpointDetails, l *logger.Logger) {
	r = r.UseEncodedPath()

	l.Println("registering routes in the order:")
	for _, ep := range sortEndpoints(maps.Keys(details)) {
		details := details[ep]
		str := fmt.Sprintf("%s %s", ep.GetMethod(), ep.GetPath())

		if details.Cors != nil {
			r.HandleFunc(string(ep.GetPath()), details.Cors(details.Handler)).Methods(http.MethodOptions)
			r.HandleFunc(string(ep.GetPath()), details.Cors(details.Handler)).Methods(ep.GetMethod())
		} else {
			r.HandleFunc(string(ep.GetPath()), details.Handler).Methods(ep.GetMethod())
		}

		l.Printf("%q -> %p\n", str, details.Handler)
	}
}

func StartServerWithEndpoints(params ServerParameters, details map[api.Endpoint]EndpointDetails) {
	l := logger.New("StartServerWithEndpoints")
	StartServer(params, func(r *mux.Router) {
		registerer(r, details, l)
	})
}
