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
	"logbook/internal/web/router/reception"
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

	r := http.NewServeMux()

	s := apicfg.Internal.Services.Registry
	agent := reception.NewAgent(deplycfg, l)
	agent.RegisterForInternal(map[api.Endpoint]reception.HandlerFunc{
		s.Endpoints.ListInstances:    e.ListInstances,
		s.Endpoints.RecheckInstance:  e.RecheckInstance,
		s.Endpoints.RegisterInstance: e.RegisterInstance,
	})
	err = agent.RegisterCommonalities()
	if err != nil {
		return fmt.Errorf("agent.RegisterCommonalities: %w", err)
	}

	err = router.StartServer(router.ServerParameters{
		Port:     deplycfg.Ports.Registry,
		Router:   deplycfg.Router,
		ServeMux: r,
		TlsCrt:   args.TlsCertificate,
		TlsKey:   args.TlsKey,
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
