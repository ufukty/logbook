package main

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/endpoints"
	"logbook/cmd/objectives/service"
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig(service.ReadConfig)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	defer internalsd.Stop()
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{})
	defer sc.Stop()

	app := app.New(pool)
	eps := endpoints.New(app)

	s := apicfg.Public.Services.Objectives
	router.StartServerWithEndpoints(router.ServerParameters{
		Address: args.PrivateNetworkIp,
		Port:    deplcfg.Ports.Objectives,
		Router:  deplcfg.Router,
		Service: models.Objectives,
		Sidecar: sc,
		TlsCrt:  args.TlsCertificate,
		TlsKey:  args.TlsKey,
	}, map[api.Endpoint]router.EndpointDetails{
		s.Endpoints.Attach:     {Handler: eps.ReattachObjective},
		s.Endpoints.Create:     {Handler: eps.CreateObjective},
		s.Endpoints.Mark:       {Handler: eps.MarkComplete},
		s.Endpoints.Placement:  {Handler: eps.GetPlacementArray},
		s.Endpoints.RockCreate: {Handler: eps.RockCreate},
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}
