package main

import (
	"fmt"
	"logbook/cmd/discovery/app"
	"logbook/cmd/discovery/cfgs"
	"logbook/cmd/discovery/endpoints"
	"logbook/config/api"
	"logbook/internal/web/router"
	"net/http"
	"os"
)

func mainerr() error {
	args, deplycfg, apicfg, err := cfgs.Read()
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	// TODO: redis connection
	a := app.New()
	ep := endpoints.New(a)

	eps := apicfg.Internal.Services.Discovery.Endpoints
	router.StartServerWithEndpoints(router.ServerParameters{
		Router:  deplycfg.Router,
		BaseUrl: "",
		TlsCrt:  args.TlsCertificate,
		TlsKey:  args.TlsKey,
	}, map[api.Endpoint]http.HandlerFunc{
		eps.List:     ep.ListInstances,
		eps.Register: ep.RegisterInstance,
	})

	return nil
}

func main() {
	if err := mainerr(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
