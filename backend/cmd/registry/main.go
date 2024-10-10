package main

import (
	"fmt"
	"log"
	"logbook/cmd/registry/app"
	"logbook/cmd/registry/endpoints"
	"logbook/config/api"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/router"
	"net/http"
)

func Main() error {
	l := logger.New("registry")

	args, deplycfg, apicfg, err := startup.Service()
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	a := app.New(deplycfg, l)
	defer a.Stop()
	e := endpoints.New(a, l)

	registerer := func(r *http.ServeMux) error {
		s := apicfg.Internal.Services.Registry

		eps := map[api.Endpoint]http.HandlerFunc{
			s.Endpoints.ListInstances:    e.ListInstances,
			s.Endpoints.RecheckInstance:  e.RecheckInstance,
			s.Endpoints.RegisterInstance: e.RegisterInstance,
		}

		for ep, handler := range eps {
			r.HandleFunc(fmt.Sprintf("%s %s", ep.GetMethod(), api.ByService(ep)), handler)
		}

		return nil
	}

	err = router.StartServer(router.ServerParameters{
		Port:        deplycfg.Ports.Registry,
		Router:      deplycfg.Router,
		TlsCrt:      args.TlsCertificate,
		TlsKey:      args.TlsKey,
		Registerers: []router.Registerer{registerer},
	}, l)
	if err != nil {
		return fmt.Errorf("router.StartServer: %w", err)
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Println(err)
	}
}
