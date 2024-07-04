package main

import (
	"fmt"
	"log"
	"logbook/cmd/gateway/cfgs"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/web/discovery"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/router"
	"logbook/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func registerForwarders(sd *discovery.ServiceDiscovery, deplcfg *deployment.Config, apicfg *api.Config, sub *mux.Router) {
	objectives, err := forwarder.NewLoadBalancedProxy(sd, models.Objectives, deplcfg.Ports.Objectives, apicfg.Public.Services.Objectives.Path)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Creating forwarder for Objectives"))
	}
	account, err := forwarder.NewLoadBalancedProxy(sd, models.Account, deplcfg.Ports.Accounts, apicfg.Public.Services.Account.Path)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Creating forwarder for Account"))
	}
	sub.PathPrefix(apicfg.Public.Services.Objectives.Path).HandlerFunc(objectives)
	sub.PathPrefix("/account").HandlerFunc(account)
}

func perform() error {
	flags, deplcfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	sd := discovery.New(models.Environment(flags.EnvMode), flags.Discovery, time.Duration(deplcfg.ServiceDiscovery.UpdatePeriod))

	router.StartServer(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf("%s%s", apicfg.Public.Path, deplcfg.Ports.Accounts),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		registerForwarders(sd, deplcfg, apicfg, r.PathPrefix("/api/v1.0.0").Subrouter())
	})

	return nil
}

func main() {
	if err := perform(); err != nil {
		log.Fatalln(err)
	}
}
