package main

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/cfgs"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/internal/web/discoveryfile"
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

	internalsd := discoveryfile.NewFileReader(flags.InternalGateway, deplcfg.ServiceDiscovery.UpdatePeriod, discoveryfile.ServiceParams{
		Port: deplcfg.Ports.Registry,
		Tls:  true,
	})
	defer internalsd.Stop()

	app := app.New(db, internalsd)
	em := endpoints.NewManager(app)

	s := apicfg.Public.Services.Objectives
	router.StartServerWithEndpoints(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf(":%d", deplcfg.Ports.Objectives),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Attach:    em.ReattachObjective,
		s.Endpoints.Create:    em.CreateObjective,
		s.Endpoints.Mark:      em.MarkComplete,
		s.Endpoints.Placement: em.GetPlacementArray,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}
