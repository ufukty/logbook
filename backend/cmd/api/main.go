package main

import (
	"fmt"
	"log"
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net/http"
)

func Main() error {
	l := logger.New("api-gateway")

	args, deplcfg, apicfg, err := startup.ApiGateway(l)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	}, l)
	// NOTE: service registry needs to be accessed through internal gateway
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{
		models.Account,
		models.Objectives,
	}, l)
	defer sc.Stop()
	defer internalsd.Stop()
	defer sc.Stop()

	var (
		accounts   = forwarder.New(sc.InstanceSource(models.Account), models.Account, api.PathFromInternet(apicfg.Public.Services.Account), l)
		objectives = forwarder.New(sc.InstanceSource(models.Objectives), models.Objectives, api.PathFromInternet(apicfg.Public.Services.Objectives), l)
	)

	registerer := func(r *http.ServeMux) error {
		r.Handle(apicfg.Public.Services.Account.Path, http.StripPrefix(apicfg.Public.Path, accounts))
		r.Handle(apicfg.Public.Services.Objectives.Path, http.StripPrefix(apicfg.Public.Path, objectives))
		return nil
	}

	err = router.StartServer(router.ServerParameters{
		Router:      deplcfg.Router,
		Port:        deplcfg.Ports.Gateway,
		TlsCrt:      args.TlsCertificate,
		TlsKey:      args.TlsKey,
		Registerers: []router.Registerer{registerer},
	}, l)
	if err != nil {
		return fmt.Errorf("router.StartServer: %w", err)
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Println(err)
	}
}
