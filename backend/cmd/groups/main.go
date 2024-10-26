package main

import (
	"context"
	"fmt"
	"logbook/cmd/groups/app"
	public "logbook/cmd/groups/endpoints/public"
	"logbook/cmd/groups/service"
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/reception"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Main() error {
	l, args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig("groups", service.ReadConfig)
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

	a := app.New(pool)
	pub := public.New(a, l)

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterEndpoints(map[api.Endpoint]http.HandlerFunc{
		apicfg.Groups.Public.Create: pub.CreateGroup,
	}, nil)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
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
