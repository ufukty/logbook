package main

import (
	"fmt"
	"log"
	"logbook/config/api"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/models"
	"net/http"
)

func Main() error {
	l := logger.New("internal-gateway")

	args, deplcfg, apicfg, err := startup.InternalGateway(l)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	registrysd := registryfile.NewFileReader(args.RegistryService, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Registry,
		Tls:  false,
	}, l)
	defer registrysd.Stop()

	account := forwarder.New(registrysd, models.Discovery, api.PathFromInternet(apicfg.Internal.Services.Account), l)
	objectives := forwarder.New(registrysd, models.Discovery, api.PathFromInternet(apicfg.Internal.Services.Objectives), l)
	registry := forwarder.New(registrysd, models.Discovery, api.PathFromInternet(apicfg.Internal.Services.Registry), l)

	router.StartServer(router.ServerParameters{
		Router: deplcfg.Router,
		Port:   deplcfg.Ports.Internal,
		TlsCrt: args.TlsCertificate,
		TlsKey: args.TlsKey,
	}, func(r *http.ServeMux) {
		// r = r.UseEncodedPath()
		r.Handle(apicfg.Internal.Services.Account.Path, http.StripPrefix(apicfg.Internal.Path, account))
		r.Handle(apicfg.Internal.Services.Objectives.Path, http.StripPrefix(apicfg.Internal.Path, objectives))
		r.Handle(apicfg.Internal.Services.Registry.Path, http.StripPrefix(apicfg.Internal.Path, registry))
	}, l)

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
