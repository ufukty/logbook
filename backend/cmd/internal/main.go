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

	"github.com/gorilla/mux"
)

func mainerr() error {
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

	accountfwd := forwarder.New(registrysd, models.Discovery, api.PathFromInternet(apicfg.Internal.Services.Account), l)
	objectivesfwd := forwarder.New(registrysd, models.Discovery, api.PathFromInternet(apicfg.Internal.Services.Objectives), l)
	registryfwd := forwarder.New(registrysd, models.Discovery, api.PathFromInternet(apicfg.Internal.Services.Registry), l)

	router.StartServer(router.ServerParameters{
		Router: deplcfg.Router,
		Port:   deplcfg.Ports.Internal,
		TlsCrt: args.TlsCertificate,
		TlsKey: args.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Internal.Path).Subrouter()
		sub.PathPrefix(apicfg.Internal.Services.Account.Path).HandlerFunc(accountfwd.Handler)
		sub.PathPrefix(apicfg.Internal.Services.Objectives.Path).HandlerFunc(objectivesfwd.Handler)
		sub.PathPrefix(apicfg.Internal.Services.Registry.Path).HandlerFunc(registryfwd.Handler)
	}, l)

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		log.Fatalln(err)
	}
}
