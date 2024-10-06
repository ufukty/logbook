package main

import (
	"fmt"
	"logbook/cmd/registry/app"
	"logbook/cmd/registry/endpoints"
	"logbook/config/api"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/router"
	"os"
)

func Main() error {
	l := logger.New("registry")

	args, deplycfg, apicfg, err := startup.Service()
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	a := app.New(deplycfg, l)
	defer a.Stop()
	eps := endpoints.New(a, l)

	s := apicfg.Internal.Services.Registry
	router.StartServerWithEndpoints(router.ServerParameters{
		Port:   deplycfg.Ports.Registry,
		Router: deplycfg.Router,
		TlsCrt: args.TlsCertificate,
		TlsKey: args.TlsKey,
	}, map[api.Endpoint]router.EndpointDetails{
		s.Endpoints.ListInstances:    {Handler: eps.ListInstances},
		s.Endpoints.RecheckInstance:  {Handler: eps.RecheckInstance},
		s.Endpoints.RegisterInstance: {Handler: eps.RegisterInstance},
	}, l)

	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
