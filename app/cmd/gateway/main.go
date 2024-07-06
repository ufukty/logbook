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

	objectives, err := forwarder.NewLoadBalancedProxy(sd, models.Objectives, deplcfg.Ports.Objectives, api.PathFromInternet(apicfg.Public.Services.Objectives))
	if err != nil {
		return fmt.Errorf("creating forwarder for objectives: %w", err)
	}
	account, err := forwarder.NewLoadBalancedProxy(sd, models.Account, deplcfg.Ports.Accounts, api.PathFromInternet(apicfg.Public.Services.Account))
	if err != nil {
		return fmt.Errorf("creating forwarder for account: %w", err)
	}

	router.StartServer(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf("%s%s", deplcfg.Api.Domain, deplcfg.Ports.Gateway),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Public.Path).Subrouter()
		sub.PathPrefix(apicfg.Public.Services.Objectives.Path).HandlerFunc(objectives)
		sub.PathPrefix(apicfg.Public.Services.Account.Path).HandlerFunc(account)
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		log.Fatalln(err)
	}
}
