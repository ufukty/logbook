package main

import (
	"context"
	"fmt"
	"log"
	objectives "logbook/cmd/objectives/client"
	profiles "logbook/cmd/profiles/client"
	"logbook/cmd/registration/app"
	"logbook/cmd/registration/endpoints"
	"logbook/cmd/registration/service"
	registry "logbook/cmd/registry/client"
	sessions "logbook/cmd/sessions/client"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/reception"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	l, args, srvcfg, deplcfg, err := startup.ServiceWithCustomConfig("registration", service.ReadConfig)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	}, l)
	defer internalsd.Stop()

	reg := registry.NewClient(balancer.NewProxied(internalsd, models.Registry))
	sc := sidecar.New(reg, deplcfg, []models.Service{models.Objectives}, l)
	defer sc.Stop()

	a := &app.App{
		Objectives: objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives))),
		Sessions:   sessions.NewClient(balancer.New(sc.InstanceSource(models.Sessions))),
		Profiles:   profiles.NewClient(balancer.New(sc.InstanceSource(models.Profiles))),
	}
	pub := endpoints.NewPublic(a, l)
	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterEndpoints(pub, nil)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
	}

	// TODO: tls between services needs certs per host(name)
	err = router.StartServer(router.ServerParameters{
		Address:  args.PrivateNetworkIp,
		Port:     deplcfg.Ports.Registration,
		Router:   deplcfg.Router,
		Service:  models.Registration,
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
	}
}
