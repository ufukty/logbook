package main

import (
	"flag"
	"fmt"
	"logbook/cmd/tasks/app"
	"logbook/cmd/tasks/database"
	"logbook/cmd/tasks/endpoints"
	"logbook/config"
	"logbook/internal/utilities/reflux"
	"logbook/internal/web/paths"
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

	// sd := serviced.New(cfg.ServiceDiscoveryConfig, cfg.ServiceDiscoveryUpdatePeriod)
	app := app.New(db)
	em := endpoints.NewManager(app)

	var handlers = map[paths.Endpoint]http.HandlerFunc{
		config.ObjectivesGetPlacementArray: em.GetPlacementArray,
		config.ObjectivesCreate:            em.CreateTask,
	}

	router.StartRouter(":"+cfg.RouterPrivate, &cfg.RouterParameters, paths.RouteRegisterer(handlers))
	router.Wait(&cfg.RouterParameters)
}
