package main

import (
	"fmt"
	"logbook/cmd/servicereg/app"
	"logbook/cmd/servicereg/cfgs"
	"logbook/cmd/servicereg/endpoints"
	"logbook/config/api"
	"logbook/internal/web/router"
	"net/http"
	"os"
	"time"
)

func mainerr() error {
	args, deplycfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	a := app.New(time.Minute, 2*time.Minute)
	defer a.Stop()

	ep := endpoints.New(a)

	eps := apicfg.Internal.Services.Discovery.Endpoints
	router.StartServerWithEndpoints(router.ServerParameters{
		Router:  deplycfg.Router,
		BaseUrl: fmt.Sprintf(":%d", deplycfg.Ports.Discovery),
		TlsCrt:  args.TlsCertificate,
		TlsKey:  args.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		eps.ListInstances:    ep.ListInstances,
		eps.RecheckInstance:  ep.RecheckInstance,
		eps.RegisterInstance: ep.RegisterInstance,
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
