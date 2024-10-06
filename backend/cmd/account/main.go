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
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/router/cors"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	l := logger.New("account")

	args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig(service.ReadConfig, l)
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
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{
		models.Objectives,
	}, l)
	defer sc.Stop()

	objectivesctl := objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives)), apicfg)

	app := app.New(pool, apicfg, objectivesctl)
	em := endpoints.New(app, l)
	s := apicfg.Public.Services.Account

	origin, err := url.JoinPath(deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}
	c := cors.Same(origin)
	// TODO: tls between services needs certs per host(name)
	router.StartServerWithEndpoints(router.ServerParameters{
		Address: args.PrivateNetworkIp,
		Port:    deplcfg.Ports.Accounts,
		Router:  deplcfg.Router,
		Service: models.Account,
		Sidecar: sc,
		TlsCrt:  args.TlsCertificate,
		TlsKey:  args.TlsKey,
	}, map[api.Endpoint]router.EndpointDetails{
		s.Endpoints.CreateUser:    {em.CreateUser, c},
		s.Endpoints.CreateProfile: {em.CreateProfile, c},
		s.Endpoints.Login:         {em.Login, c},
		s.Endpoints.Logout:        {em.Logout, c},
		s.Endpoints.Whoami:        {em.WhoAmI, c},
	}, l)

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
