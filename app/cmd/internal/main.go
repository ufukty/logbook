package main

import (
	"fmt"
	"log"
	"logbook/cmd/gateway/cfgs"
	"logbook/internal/web/discoveryfile"
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

	sd := discoveryfile.NewFileReader(flags.Discovery, deplcfg.ServiceDiscovery.UpdatePeriod, discoveryfile.ServiceParams{
		Port: deplcfg.Ports.Discovery,
		Tls:  false,
	})
	defer sd.Stop()

	discovery := forwarder.New(sd, models.Discovery, apicfg.Internal.Services.Registry.Path)

	router.StartServer(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf(":%d", deplcfg.Ports.Internal),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Public.Path).Subrouter()
		sub.PathPrefix(apicfg.Internal.Services.Registry.Path).HandlerFunc(discovery.Handler)
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		log.Fatalln(err)
	}
}
