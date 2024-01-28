package main

import (
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

func main() {
	godotenv.Load("../.env")
	godotenv.Load("../.local.env")
	db, err := database.New(os.Getenv("DSN"))
	if err != nil {
		panic(fmt.Errorf("creating database instance: %w", err))
	}
	defer db.Close()

	cfg := config.Read()
	reflux.Print(cfg.Tasks)
	// sd := serviced.New(cfg.Tasks.ServiceDiscoveryConfig, cfg.Tasks.ServiceDiscoveryUpdatePeriod)
	app := app.New(db)
	em := endpoints.NewManager(app)

	var handlers = map[paths.Endpoint]http.HandlerFunc{
		config.ObjectivesGetPlacementArray: em.GetPlacementArray,
		config.ObjectivesCreate:            em.CreateTask,
	}

	router.StartRouter(":"+cfg.Tasks.RouterPrivate, &cfg.Tasks.RouterParameters, paths.RouteRegisterer(handlers))
	router.Wait(&cfg.Tasks.RouterParameters)
}
