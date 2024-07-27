package main

import (
	"fmt"
	"logbook/cmd/registry/app"
	"logbook/cmd/registry/cfgs"
	"logbook/cmd/registry/endpoints"
	"logbook/config/api"
	"logbook/internal/web/router"
	"net/http"
	"os"
	"time"
)

func Main() error {
	args, deplycfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	a := app.New(time.Minute, 2*time.Minute)
	defer a.Stop()
	eps := endpoints.New(a)

	s := apicfg.Internal.Services.Registry
	router.StartServerWithEndpoints(router.ServerParameters{
		Router: deplycfg.Router,
		Port:   deplycfg.Ports.Registry,
		TlsCrt: args.TlsCertificate,
		TlsKey: args.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.ListInstances:    eps.ListInstances,
		s.Endpoints.RecheckInstance:  eps.RecheckInstance,
		s.Endpoints.RegisterInstance: eps.RegisterInstance,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
