package main

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/api/public"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/service"
	registry "logbook/cmd/registry/client"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	l := logger.New("objectives")

	args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig(service.ReadConfig, l)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()

	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	}, l)
	defer internalsd.Stop()
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{}, l)
	defer sc.Stop()

	app := app.New(pool, l)

	pub := public.New(apicfg, deplcfg, app, sc, l)
	pri := public.New(apicfg, deplcfg, app, sc, l)

	router.StartServer(router.ServerParameters{
		Address: args.PrivateNetworkIp,
		Port:    deplcfg.Ports.Objectives,
		Router:  deplcfg.Router,
		Service: models.Objectives,
		Sidecar: sc,
		TlsCrt:  args.TlsCertificate,
		TlsKey:  args.TlsKey,
	}, func(r *http.ServeMux) {
		pub.Register(r)
		pri.Register(r)
	}, l)

	return nil
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}
