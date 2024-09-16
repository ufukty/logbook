package main

import (
	"context"
	"fmt"
	"log"
	registry "logbook/cmd/registry/client"
	"logbook/cmd/tags/app"
	"logbook/cmd/tags/endpoints"
	"logbook/cmd/tags/service"
	"logbook/config/api"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig(service.ReadConfig)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg.ServiceDiscovery.UpdatePeriod, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	defer internalsd.Stop()
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg.ServiceDiscovery.UpdatePeriod, []models.Service{})
	defer sc.Stop()

	app := app.New(pool, apicfg, internalsd)
	eps := endpoints.New(app)

	// TODO: tls between services needs certs per host(name)
	s := apicfg.Public.Services.Tags
	router.StartServerWithEndpoints(router.ServerParameters{
		Address: args.PrivateNetworkIp,
		Port:    deplcfg.Ports.Tags,
		Router:  deplcfg.Router,
		Service: models.Tags,
		Sidecar: sc,
		TlsCrt:  args.TlsCertificate,
		TlsKey:  args.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Assign:   eps.TagAssign,
		s.Endpoints.Creation: eps.TagCreation,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
