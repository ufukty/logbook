package main

import (
	"fmt"
	"log"
	discovery "logbook/cmd/discovery/client"
	"logbook/cmd/gateway/cfgs"
	"logbook/config/api"
	"logbook/internal/web/balancer"
	"logbook/internal/web/discoveryfile"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/router"
	"logbook/models"
	"time"

	"github.com/gorilla/mux"
)

func mainerr() error {
	flags, deplcfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	is := discoveryfile.NewFileReader(flags.Discovery, time.Second, discoveryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	defer is.Stop()
	lb := balancer.New(is)
	dsctl := discovery.NewClient(lb)
	
	// FIXME:
	dsctl.ListInstances()
	
	// FIXME:
	objectives, err := forwarder.New(, models.Objectives, api.PathFromInternet(apicfg.Public.Services.Objectives))
	if err != nil {
		return fmt.Errorf("creating forwarder for objectives: %w", err)
	}
	
	// FIXME:
	account, err := forwarder.New(, models.Account, api.PathFromInternet(apicfg.Public.Services.Account))
	if err != nil {
		return fmt.Errorf("creating forwarder for account: %w", err)
	}

	router.StartServer(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf(":%d", deplcfg.Ports.Gateway),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Public.Path).Subrouter()
		sub.PathPrefix(apicfg.Public.Services.Objectives.Path).HandlerFunc(objectives.Handler)
		sub.PathPrefix(apicfg.Public.Services.Account.Path).HandlerFunc(account.Handler)
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		log.Fatalln(err)
	}
}
