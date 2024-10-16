package main

import (
	"fmt"
	"log"
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/internal/logger"
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
	l := logger.New("api-gateway")

	args, deplcfg, apicfg, err := startup.ApiGateway(l)
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

	s := apicfg.Public.Services

	var (
		accounts   = forwarder.New(sc.InstanceSource(models.Account), models.Account, api.ByGateway(s.Account), l)
		objectives = forwarder.New(sc.InstanceSource(models.Objectives), models.Objectives, api.ByGateway(s.Objectives), l)
	)

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterForwarders(apicfg.Public.Path, map[api.Addressable]*forwarder.LoadBalancedReverseProxy{
		s.Account:    accounts,
		s.Objectives: objectives,
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForwarders: %w", err)
	}
	err = agent.RegisterCommonalities()
	if err != nil {
		return fmt.Errorf("agent.RegisterCommonalities: %w", err)
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
