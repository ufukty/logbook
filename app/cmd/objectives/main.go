package main

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/args"
	"logbook/internal/utilities/reflux"
	"logbook/internal/web/router"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func Main() error {
	flags, err := args.Parse()
	if err != nil {
		return fmt.Errorf("parsing args: %w", err)
	}

	godotenv.Load(flags.Environment)
	db, err := database.New(os.Getenv("DSN"))
	if err != nil {
		return fmt.Errorf("creating database instance: %w", err)
	}
	defer db.Close()

	cfg := deployment.Read(flags.Config).Tasks
	reflux.Print(cfg)

	apicfg, err := api.ReadConfig("../../api.yml")
	if err != nil {
		return fmt.Errorf("reading api config: %w", err)
	}

	// sd := serviced.New(cfg.ServiceDiscoveryConfig, cfg.ServiceDiscoveryUpdatePeriod)
	app := app.New(db)
	em := endpoints.NewManager(app)

	eps := apicfg.Gateways.Public.Services.Objectives.Endpoints
	router.StartServer(":"+cfg.RouterPrivate, false, cfg.RouterParameters, map[api.Endpoint]http.HandlerFunc{
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
