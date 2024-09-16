package main

import (
	"context"
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/cmd/account/endpoints"
	"logbook/cmd/account/service"
	objectives "logbook/cmd/objectives/client"
	registry "logbook/cmd/registry/client"
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
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{
		models.Objectives,
	})
	defer sc.Stop()

	objectivesctl := objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives)), apicfg)

	app := app.New(pool, apicfg, objectivesctl)
	em := endpoints.New(app)
	s := apicfg.Public.Services.Account

	// TODO: tls between services needs certs per host(name)
	router.StartServerWithEndpoints(router.ServerParameters{
		Address: args.PrivateNetworkIp,
		Port:    deplcfg.Ports.Accounts,
		Router:  deplcfg.Router,
		Service: models.Account,
		Sidecar: sc,
		TlsCrt:  args.TlsCertificate,
		TlsKey:  args.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Create:        em.CreateUser,
		s.Endpoints.CreateProfile: em.CreateProfile,
		s.Endpoints.Login:         em.Login,
		s.Endpoints.Logout:        em.Logout,
		s.Endpoints.Whoami:        em.WhoAmI,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
