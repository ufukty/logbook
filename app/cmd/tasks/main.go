package main

import (
	"flag"
	"fmt"
	"logbook/cmd/tasks/app"
	"logbook/cmd/tasks/database"
	"logbook/cmd/tasks/endpoints"
	"logbook/config"
	"logbook/internal/utilities/reflux"
	"logbook/internal/web/api"
	"logbook/internal/web/router"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getConfigPath() string {
	var configpath string
	flag.StringVar(&configpath, "config", "", "")
	flag.Parse()
	return configpath
}

func main() {
	godotenv.Load("../.env")
	godotenv.Load("../.local.env")

	db, err := database.New(os.Getenv("DSN"))
	if err != nil {
		panic(fmt.Errorf("creating database instance: %w", err))
	}
	defer db.Close()

	cfg := config.Read(getConfigPath()).Tasks
	reflux.Print(cfg)

	apicfg, err := api.ReadConfig("../../api.yml")
	if err != nil {
		panic(fmt.Errorf("reading api config: %w", err))
	}

	// sd := serviced.New(cfg.ServiceDiscoveryConfig, cfg.ServiceDiscoveryUpdatePeriod)
	app := app.New(db)
	em := endpoints.NewManager(app)

	eps := apicfg.Gateways.Public.Services.Objectives.Endpoints
	var handlers = map[api.Endpoint]http.HandlerFunc{
		eps.Attach:    em.ReattachObjective,
		eps.Create:    em.CreateTask,
		eps.Mark:      em.MarkComplete,
		eps.Placement: em.GetPlacementArray,
	}

	router.StartServer(":"+cfg.RouterPrivate, false, cfg.RouterParameters, api.RouteRegisterer(handlers))
}
