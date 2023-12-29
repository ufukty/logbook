package main

import (
	"fmt"
	"logbook/cmd/tasks/database"
	"logbook/cmd/tasks/endpoints"
	"logbook/config"
	"logbook/config/reader"
	"logbook/internal/web/paths"
	"logbook/internal/web/router"
	"net/http"
)

func main() {
	db, err := database.New("postgres://ufuktan:password@localhost:5432/logbook_dev")
	if err != nil {
		panic(fmt.Errorf("creating database instance: %w", err))
	}
	defer db.Close()

	cfg := reader.GetConfig()
	// sd := serviced.New(cfg.Tasks.ServiceDiscoveryConfig, cfg.Tasks.ServiceDiscoveryUpdatePeriod)
	em := endpoints.NewManager(db)

	reader.Print(cfg.Tasks)

	var handlers = map[paths.Endpoint]http.HandlerFunc{
		config.TaskList: em.ListTasks,
	}

	router.StartRouter(":"+cfg.Tasks.RouterPrivate, &cfg.Tasks.RouterParameters, paths.RouteRegisterer(handlers))
	router.Wait(&cfg.Tasks.RouterParameters)
}
