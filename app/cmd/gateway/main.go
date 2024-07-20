package main

import (
	"fmt"
	"log"
	"logbook/cmd/gateway/cfgs"
	"logbook/cmd/gateway/forwarders"
	"logbook/internal/web/router"

	"github.com/gorilla/mux"
)

func mainerr() error {
	flags, deplcfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	fws, err := forwarders.New(flags, deplcfg, apicfg)
	defer fws.Stop()

	router.StartServer(router.ServerParameters{
		Router:  deplcfg.Router,
		BaseUrl: fmt.Sprintf(":%d", deplcfg.Ports.Gateway),
		TlsCrt:  flags.TlsCertificate,
		TlsKey:  flags.TlsKey,
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
