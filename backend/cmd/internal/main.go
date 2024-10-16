package main

import (
	"fmt"
	"log"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/reception"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
)

func Main() error {
	l := logger.New("internal-gateway")

	args, deplcfg, apicfg, err := startup.InternalGateway(l)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	registrysd := registryfile.NewFileReader(args.RegistryService, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Registry,
		Tls:  false,
	}, l)
	defer registrysd.Stop()

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterForwarders(map[string]*forwarder.LoadBalancedReverseProxy{
		apicfg.InternalGateway.Services.Registry: forwarder.New(registrysd, l),
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForwarders: %w", err)
	}
	// err = agent.RegisterCommonalities()
	// if err != nil {
	// 	return fmt.Errorf("agent.RegisterCommonalities: %w", err)
	// }

	router.StartServer(router.ServerParameters{
		Router:   deplcfg.Router,
		Port:     deplcfg.Ports.Internal,
		ServeMux: agent.Mux(),
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
