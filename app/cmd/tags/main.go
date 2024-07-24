package main

import (
	"fmt"
	"log"
	"logbook/cmd/tags/app"
	"logbook/cmd/tags/cfgs"
	"logbook/cmd/tags/database"
	"logbook/cmd/tags/endpoints"
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
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	defer internalsd.Stop()

	app := app.New(db, apicfg, internalsd)
	eps := endpoints.New(app)

	// TODO: tls between services needs certs per host(name)
	s := apicfg.Public.Services.Tags
	router.StartServerWithEndpoints(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf(":%d", deplcfg.Ports.Tags),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
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
