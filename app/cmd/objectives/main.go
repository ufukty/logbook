package main

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/endpoints"
	"logbook/cmd/objectives/service"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/args"
	"logbook/internal/utilities/reflux"
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
	reflux.Print(srvcfg)

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
	em := endpoints.NewManager(app)

	eps := apicfg.Gateways.Public.Services.Objectives.Endpoints
	router.StartServer(router.ServerParameters{
		BaseUrl:        depl.Ports.Objectives,
		Tls:            false,
		RequestTimeout: depl.Router.RequestTimeout,
	}, map[api.Endpoint]http.HandlerFunc{
		eps.Attach:    em.ReattachObjective,
		eps.Create:    em.CreateTask,
		eps.Mark:      em.MarkComplete,
		eps.Placement: em.GetPlacementArray,
	})

	return nil
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}
