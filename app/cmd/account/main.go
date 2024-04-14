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
	"logbook/internal/utilities/reflux"
	"logbook/internal/web/logger"
	"logbook/internal/web/router"
	"net/http"
)

func readConfigs() (*service.Config, *deployment.Config, *api.Config, error) {
	flags, err := args.Parse()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}
	l := logger.NewLogger("readConfigs")

	srvcfg, err := service.ReadConfig(flags.Service)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading service config: %w", err)
	}
	l.Println("service config:")
	reflux.Print(srvcfg)

	deplcfg, err := deployment.ReadConfig(flags.Deployment)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(flags.Api)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return &srvcfg, &deplcfg, &apicfg, nil
}

func Main() error {
	srvcfg, deplcfg, apicfg, err := readConfigs()
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	db, err := database.New(srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("creating database instance: %w", err)
	}
	defer db.Close()

	// sd := serviced.New(cfg.ServiceDiscoveryConfig, cfg.ServiceDiscoveryUpdatePeriod)
	app := app.New(db)
	em := endpoints.New(app)

	eps := apicfg.Gateways.Public.Services.Account.Endpoints
	router.StartServer(router.ServerParameters{
		BaseUrl:        deplcfg.Ports.Accounts,
		Tls:            false,
		RequestTimeout: deplcfg.Router.RequestTimeout,
	}, map[api.Endpoint]http.HandlerFunc{
		eps.Create:        em.CreateUser,
		eps.CreateProfile: em.CreateProfile,
		eps.CreateSession: em.CreateSession,
		eps.Whoami:        em.WhoAmI,
		eps.Logout:        em.Logout,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
