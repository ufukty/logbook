package main

import (
	"fmt"
	"log"
	registry "logbook/cmd/registry/client"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/reception"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"
)

func Main() error {
	l, args, deplcfg, err := startup.ApiGateway("api-gateway")
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	}, l)
	defer internalsd.Stop()

	// NOTE: service registry needs to be accessed through internal gateway
	reg := registry.NewClient(balancer.NewProxied(internalsd, models.Registry))
	sc := sidecar.New(reg, deplcfg, []models.Service{
		models.Objectives,
		models.Profiles,
		models.Registration,
		models.Users,
	}, l)
	defer sc.Stop()

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterForwarders(map[models.Service]*forwarder.LoadBalancedReverseProxy{
		models.Users:        forwarder.New(sc.InstanceSource(models.Users), deplcfg, l),
		models.Objectives:   forwarder.New(sc.InstanceSource(models.Objectives), deplcfg, l),
		models.Profiles:     forwarder.New(sc.InstanceSource(models.Profiles), deplcfg, l),
		models.Registration: forwarder.New(sc.InstanceSource(models.Registration), deplcfg, l),
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForwarders: %w", err)
	}

	err = router.StartServer(router.ServerParameters{
		Port:     deplcfg.Ports.Gateway,
		Router:   deplcfg.Router,
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
