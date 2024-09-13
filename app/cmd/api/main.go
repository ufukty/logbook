package main

import (
	"fmt"
	"log"
	"logbook/cmd/api/forwarders"
	"logbook/internal/startup"
	"logbook/internal/web/router"

	"github.com/gorilla/mux"
)

func mainerr() error {
	args, deplcfg, apicfg, err := startup.EverythingForApiGateway()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	fws, err := forwarders.New(args, deplcfg, apicfg)
	if err != nil {
		return fmt.Errorf("forwarders.New: %w", err)
	}
	defer fws.Stop()

	router.StartServer(router.ServerParameters{
		Router: deplcfg.Router,
		Port:   deplcfg.Ports.Gateway,
		TlsCrt: args.TlsCertificate,
		TlsKey: args.TlsKey,
	}, func(r *mux.Router) {
		r = r.UseEncodedPath()
		sub := r.PathPrefix(apicfg.Public.Path).Subrouter()
		sub.PathPrefix(apicfg.Public.Services.Account.Path).HandlerFunc(fws.Accounts.Handler)
		sub.PathPrefix(apicfg.Public.Services.Objectives.Path).HandlerFunc(fws.Objectives.Handler)
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		log.Fatalln(err)
	}
}
