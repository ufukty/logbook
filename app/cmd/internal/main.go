package main

import (
	"fmt"
	"log"
	"logbook/cmd/gateway/cfgs"
	"logbook/config/api"
	"logbook/internal/web/discovery"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/router"
	"logbook/models"

	"github.com/gorilla/mux"
)

func mainerr() error {
	flags, deplcfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	sd := discovery.New(models.Environment(flags.EnvMode), flags.Discovery, deplcfg.ServiceDiscovery.UpdatePeriod)

	discovery, err := forwarder.New(sd, models.Discovery, deplcfg.Ports.Discovery, api.PathFromInternet(apicfg.Internal.Services.Discovery))
	if err != nil {
		return fmt.Errorf("creating forwarder for objectives: %w", err)
	}

	router.StartServer(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: deplcfg.Ports.Internal,
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Public.Path).Subrouter()
		sub.PathPrefix(apicfg.Internal.Services.Discovery.Path).HandlerFunc(discovery.Handler)
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		log.Fatalln(err)
	}
}
