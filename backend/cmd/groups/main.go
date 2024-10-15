package main

import (
	"context"
	"fmt"
	"logbook/cmd/groups/api/public"
	"logbook/cmd/groups/service"
	registry "logbook/cmd/registry/client"
	"logbook/internal/logger"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/router/reception"
	"logbook/internal/web/sidecar"
	"logbook/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	l := logger.New("groups")

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

	pub := public.New(apicfg, deplcfg, pool, sc, l)

	agent := reception.NewAgent(deplcfg, l)
	err = pub.Register(agent)
	if err != nil {
		return fmt.Errorf("pub.Register: %w", err)
	}

	router.StartServer(router.ServerParameters{
		Address:  args.PrivateNetworkIp,
		Port:     deplcfg.Ports.Objectives,
		Router:   deplcfg.Router,
		ServeMux: agent.Mux(),
		Service:  models.Groups,
		Sidecar:  sc,
		TlsCrt:   args.TlsCertificate,
		TlsKey:   args.TlsKey,
	}, l)

	return nil
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}
