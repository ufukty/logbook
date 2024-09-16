package main

import (
	"fmt"
	"log"
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"

	"github.com/gorilla/mux"
)

func Main() error {
	args, deplcfg, apicfg, err := startup.ApiGateway()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg.ServiceDiscovery.UpdatePeriod, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	// NOTE: service registry needs to be accessed through internal gateway
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{
		models.Account,
		models.Objectives,
	})
	defer sc.Stop()
	defer internalsd.Stop()
	defer sc.Stop()

	accounts := forwarder.New(sc.InstanceSource(models.Account), models.Account, api.PathFromInternet(apicfg.Public.Services.Account))
	objectives := forwarder.New(sc.InstanceSource(models.Objectives), models.Objectives, api.PathFromInternet(apicfg.Public.Services.Objectives))

	router.StartServer(router.ServerParameters{
		Router: deplcfg.Router,
		Port:   deplcfg.Ports.Gateway,
		TlsCrt: args.TlsCertificate,
		TlsKey: args.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Public.Path).Subrouter()
		sub.PathPrefix(apicfg.Public.Services.Account.Path).HandlerFunc(accounts.Handler)
		sub.PathPrefix(apicfg.Public.Services.Objectives.Path).HandlerFunc(objectives.Handler)
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
