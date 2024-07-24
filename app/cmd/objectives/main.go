package main

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/cfgs"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"net/http"
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
		Port: deplcfg.Ports.Registry,
		Tls:  true,
	})
	defer internalsd.Stop()

	app := app.New(db, internalsd)
	eps := endpoints.New(app)

	s := apicfg.Public.Services.Objectives
	router.StartServerWithEndpoints(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf(":%d", deplcfg.Ports.Objectives),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Attach:    eps.ReattachObjective,
		s.Endpoints.Create:    eps.CreateObjective,
		s.Endpoints.Mark:      eps.MarkComplete,
		s.Endpoints.Placement: eps.GetPlacementArray,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}
