package main

import (
	"context"
	"fmt"
	"log"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/endpoints"
	"logbook/cmd/objectives/service"
	registry "logbook/cmd/registry/client"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/reception"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	l, args, srvcfg, deplcfg, _, err := startup.ServiceWithCustomConfig("objectives", service.ReadConfig)
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

	reg := registry.NewClient(balancer.NewProxied(internalsd, models.Registry))
	sc := sidecar.New(reg, deplcfg, []models.Service{}, l)
	defer sc.Stop()

	a := app.New(pool, l)
	pub := endpoints.NewPublic(a, l)
	priv := endpoints.NewPrivate(a, l)

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterEndpoints(pub, priv)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
	}

	err = router.StartServer(router.ServerParameters{
		Address:  args.PrivateNetworkIp,
		Port:     deplcfg.Ports.Objectives,
		Router:   deplcfg.Router,
		ServeMux: agent.Mux(),
		Service:  models.Objectives,
		Sidecar:  sc,
		TlsCrt:   args.TlsCertificate,
		TlsKey:   args.TlsKey,
	}, l)
	if err != nil {
		return fmt.Errorf("router.StartServer: %w", err)
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Println(err)
	}
}
