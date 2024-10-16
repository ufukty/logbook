package main

import (
	"fmt"
	"log"
	"logbook/config/api"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/reception"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/models"
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

	s := apicfg.Internal.Services

	var (
		account    = forwarder.New(registrysd, models.Discovery, api.ByGateway(s.Account), l)
		objectives = forwarder.New(registrysd, models.Discovery, api.ByGateway(s.Objectives), l)
		registry   = forwarder.New(registrysd, models.Discovery, api.ByGateway(s.Registry), l)
	)

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterForwarders(apicfg.Internal.Path, map[api.Addressable]*forwarder.LoadBalancedReverseProxy{
		s.Account:    account,
		s.Objectives: objectives,
		s.Registry:   registry,
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForwarders: %w", err)
	}
	err = agent.RegisterCommonalities()
	if err != nil {
		return fmt.Errorf("agent.RegisterCommonalities: %w", err)
	}

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
