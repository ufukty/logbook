package main

import (
	"fmt"
	"log"
	"logbook/config/api"
	"logbook/internal/startup"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/models"

	"github.com/gorilla/mux"
)

func mainerr() error {
	args, deplcfg, apicfg, err := startup.EverythingForInternalGateway()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	registrysd := registryfile.NewFileReader(args.RegistryService, deplcfg.ServiceDiscovery.UpdatePeriod, registryfile.ServiceParams{
		Port: deplcfg.Ports.Registry,
		Tls:  false,
	})
	defer registrysd.Stop()

	registryfwd := forwarder.New(registrysd, models.Discovery, api.PathFromInternet(apicfg.Internal.Services.Registry))

	router.StartServer(router.ServerParameters{
		Router: deplcfg.Router,
		Port:   deplcfg.Ports.Internal,
		TlsCrt: args.TlsCertificate,
		TlsKey: args.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Internal.Path).Subrouter()
		sub.PathPrefix(apicfg.Internal.Services.Registry.Path).HandlerFunc(registryfwd.Handler)
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		log.Fatalln(err)
	}
}
