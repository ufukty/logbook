package main

import (
	"fmt"
	"log"

	"logbook/internal/startup"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/register"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/models"
)

func Main() error {
	l, args, deplcfg, err := startup.InternalGateway()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	registrysd := registryfile.NewFileReader(args.RegistryService, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Registry,
		Tls:  false,
	}, l)
	defer registrysd.Stop()

	r, err := register.RegisterForwarders(deplcfg, l, map[models.Service]*forwarder.LoadBalancedReverseProxy{
		models.Registry: forwarder.New(registrysd, deplcfg, l),
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForwarders: %w", err)
	}

	router.StartServer(router.ServerParameters{
		Router:   deplcfg.Router,
		Port:     deplcfg.Ports.Internal,
		ServeMux: r,
		TlsCrt:   args.TlsCertificate,
		TlsKey:   args.TlsKey,
	}, l)

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
