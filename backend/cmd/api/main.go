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
	l, args, deplcfg, apicfg, err := startup.ApiGateway("api-gateway")
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	}, l)
	defer internalsd.Stop()

	// NOTE: service registry needs to be accessed through internal gateway
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{
		models.Account,
		models.Objectives,
	}, l)
	defer sc.Stop()

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterForwarders(map[string]*forwarder.LoadBalancedReverseProxy{
		apicfg.ApiGateway.Services.Account:    forwarder.New(sc.InstanceSource(models.Account), deplcfg, l),
		apicfg.ApiGateway.Services.Objectives: forwarder.New(sc.InstanceSource(models.Objectives), deplcfg, l),
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForwarders: %w", err)
	}
	// err = agent.RegisterCommonalities()
	// if err != nil {
	// 	return fmt.Errorf("agent.RegisterCommonalities: %w", err)
	// }

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
