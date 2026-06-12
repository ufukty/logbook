package main

import (
	"fmt"
	"log"

	"logbook/cmd/registry/app"
	"logbook/cmd/registry/endpoints"
	"logbook/internal/startup"
	"logbook/internal/web/register"
	"logbook/internal/web/router"
)

func Main() error {
	l, args, deplcfg, err := startup.Service("registry")
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	a := app.New(deplcfg, l)
	defer a.Stop()
	e := endpoints.New(a, l)

	r, err := register.RegisterEndpoints(deplcfg, l, nil, e)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
	}

	err = router.StartServer(router.ServerParameters{
		Port:     deplcfg.Ports.Registry,
		Router:   deplcfg.Router,
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
