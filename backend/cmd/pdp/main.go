package main

import (
	"fmt"
	"log"
	groups "logbook/cmd/groups/client"
	objectives "logbook/cmd/objectives/client"
	"logbook/cmd/pdp/decider"
	"logbook/cmd/pdp/endpoints"
	registry "logbook/cmd/registry/client"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/reception"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"os"
)

func Main() error {
	l, args, deplcfg, _, err := startup.Service("pdp")
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	}, l)
	defer internalsd.Stop()

	reg := registry.NewClient(balancer.NewProxied(internalsd, models.Registry))
	speeddial := []models.Service{
		models.Groups,
		models.Objectives,
	}
	sc := sidecar.New(reg, deplcfg, speeddial, l)
	defer sc.Stop()

	d := decider.New(
		groups.NewClient(balancer.New(sc.InstanceSource(models.Groups))),
		objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives))),
	)
	eps := endpoints.NewPrivate(d, l)
	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterEndpoints(nil, eps)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
	}

	// TODO: tls between services needs certs per host(name)
	err = router.StartServer(router.ServerParameters{
		Address:  args.PrivateNetworkIp,
		Port:     deplcfg.Ports.Accounts,
		Router:   deplcfg.Router,
		Service:  models.Pdp,
		ServeMux: agent.Mux(),
		Sidecar:  sc,
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
		os.Exit(1)
	}
}
