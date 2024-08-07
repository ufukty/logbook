package main

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/cfgs"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/endpoints"
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net/http"
	"time"
)

func Main() error {
	flags, srvcfg, deplcfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	db, err := database.New(srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("creating database instance: %w", err)
	}
	defer db.Close()

	internalsd := registryfile.NewFileReader(flags.InternalGateway, deplcfg.ServiceDiscovery.UpdatePeriod, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	defer internalsd.Stop()
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), time.Second, []models.Service{})
	defer sc.Stop()

	app := app.New(db)
	eps := endpoints.New(app)

	s := apicfg.Public.Services.Objectives
	router.StartServerWithEndpoints(router.ServerParameters{
		Router:  deplcfg.Router,
		Address: flags.PrivateNetworkIp,
		Port:    deplcfg.Ports.Objectives,
		Sidecar: sc,
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Attach:     eps.ReattachObjective,
		s.Endpoints.Create:     eps.CreateObjective,
		s.Endpoints.Mark:       eps.MarkComplete,
		s.Endpoints.Placement:  eps.GetPlacementArray,
		s.Endpoints.RockCreate: eps.RockCreate,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}
