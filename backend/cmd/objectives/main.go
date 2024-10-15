package main

import (
	"context"
	"fmt"
	"log"
	private "logbook/cmd/objectives/api/private/endpoints"
	public "logbook/cmd/objectives/api/public/endpoints"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/service"
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
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

	a := app.New(pool, l)
	pub := public.New(a, l)
	priv := private.New(a, l)

	agent := reception.NewAgent(deplcfg, l)

	ps := apicfg.Public.Services.Objectives
	err = agent.RegisterForPublic(map[api.Endpoint]reception.HandlerFunc{
		ps.Endpoints.Attach:    pub.ReattachObjective,
		ps.Endpoints.Create:    pub.CreateObjective,
		ps.Endpoints.Mark:      pub.MarkComplete,
		ps.Endpoints.Placement: pub.GetPlacementArray,
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForPublic: %w", err)
	}

	is := apicfg.Internal.Services.Objectives
	err = agent.RegisterForInternal(map[api.Endpoint]reception.HandlerFunc{
		is.Endpoints.RockCreate: priv.RockCreate,
	})
	if err != nil {
		return fmt.Errorf("agent.RegisterForInternal: %w", err)
	}
	err = agent.RegisterCommonalities()
	if err != nil {
		return fmt.Errorf("agent.RegisterCommonalities: %w", err)
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
