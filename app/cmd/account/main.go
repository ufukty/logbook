package main

import (
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/endpoints"
	"logbook/cmd/account/service"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/args"
	"logbook/internal/web/router"
	"net/http"
)

func Main() error {
	flags, err := args.Parse()
	if err != nil {
		return fmt.Errorf("parsing args: %w", err)
	}

	srvcfg, err := service.ReadConfig(flags.Service)
	if err != nil {
		return fmt.Errorf("reading service config: %w", err)
	}

	db, err := database.New(srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("creating database instance: %w", err)
	}
	defer db.Close()

	apicfg, err := api.ReadConfig(flags.Api)
	if err != nil {
		return fmt.Errorf("reading api config: %w", err)
	}

	depl, err := deployment.ReadConfig(flags.Deployment)
	if err != nil {
		return fmt.Errorf("reading deployment environment config: %w", err)
	}

	// sd := serviced.New(cfg.ServiceDiscoveryConfig, cfg.ServiceDiscoveryUpdatePeriod)
	app := app.New(db)
	em := endpoints.New(app)

	eps := apicfg.Gateways.Public.Services.Account.Endpoints
	router.StartServer(router.ServerParameters{
		BaseUrl:        depl.Ports.Accounts,
		Tls:            false,
		RequestTimeout: depl.Router.RequestTimeout,
	}, map[api.Endpoint]http.HandlerFunc{
		eps.Create:        em.CreateUser,
		eps.CreateSession: em.CreateSession,
		eps.Whoami:        em.WhoAmI,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
