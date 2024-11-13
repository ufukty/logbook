package main

import (
	"context"
	"fmt"
	"log"
	"logbook/cmd/profiles/app"
	"logbook/cmd/profiles/endpoints"
	"logbook/cmd/profiles/service"
	registry "logbook/cmd/registry/client"
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
	l, args, srvcfg, deplcfg, err := startup.ServiceWithCustomConfig("profiles", service.ReadConfig)
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

	a := app.New(pool)

	pri := endpoints.NewPrivate(a, l)

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterEndpoints(nil, pri)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
	}

	// TODO: tls between services needs certs per host(name)
	err = router.StartServer(router.ServerParameters{
		Address:  args.PrivateNetworkIp,
		Port:     deplcfg.Ports.Profiles,
		Router:   deplcfg.Router,
		Service:  models.Profiles,
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
