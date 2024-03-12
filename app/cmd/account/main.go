package main

import (
	"flag"
	"fmt"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/endpoints"
	"logbook/config"
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

	apicfg, err := api.ReadConfig("../../api.yml")
	if err != nil {
		panic(fmt.Errorf("reading api config: %w", err))
	}

	// sd := serviced.New(cfg.ServiceDiscoveryConfig, cfg.ServiceDiscoveryUpdatePeriod)
	app := app.New(db)
	em := endpoints.New(app)

	eps := apicfg.Gateways.Public.Services.Account.Endpoints
	router.StartServer(":"+cfg.RouterPrivate, false, cfg.RouterParameters, map[api.Endpoint]http.HandlerFunc{
		eps.Create: em.CreateUser,
	})
}
