package main

import (
	"fmt"
	"log"
	"logbook/cmd/registry/app"
	"logbook/cmd/registry/endpoints"
	"logbook/internal/startup"
	"logbook/internal/web/reception"
	"logbook/internal/web/router"
)

func Main() error {
	l, args, deplycfg, err := startup.Service("registry")
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	a := app.New(deplycfg, l)
	defer a.Stop()
	e := endpoints.New(a, l)

	agent := reception.NewAgent(deplycfg, l)
	err = agent.RegisterEndpoints(nil, e)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
	}

	err = router.StartServer(router.ServerParameters{
		Port:     deplycfg.Ports.Registry,
		Router:   deplycfg.Router,
		ServeMux: agent.Mux(),
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
