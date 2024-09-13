package main

import (
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
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
	"time"
)

func Main() error {
	args, srvcfg, deplcfg, apicfg, err := startup.EverythingForServiceWithCustomServiceConfig(service.ReadConfig)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	db, err := database.New(srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("creating database instance: %w", err)
	}
	defer db.Close()

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg.ServiceDiscovery.UpdatePeriod, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	defer internalsd.Stop()
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), time.Second, []models.Service{
		models.Objectives,
	})
	defer sc.Stop()

	objectivesctl := objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives)), apicfg)

	app := app.New(db, apicfg, objectivesctl)
	em := endpoints.New(app)
	s := apicfg.Public.Services.Account

	// TODO: tls between services needs certs per host(name)
	router.StartServerWithEndpoints(router.ServerParameters{
		Router:  deplcfg.Router,
		Address: args.PrivateNetworkIp,
		Port:    deplcfg.Ports.Accounts,
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
