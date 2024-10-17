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
	l, args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig("objectives", service.ReadConfig)
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

	err = agent.RegisterEndpoints(
		map[api.Endpoint]http.HandlerFunc{
			apicfg.Objectives.Public.Attach:    pub.ReattachObjective,
			apicfg.Objectives.Public.Create:    pub.CreateObjective,
			apicfg.Objectives.Public.Mark:      pub.MarkComplete,
			apicfg.Objectives.Public.Placement: pub.GetPlacementArray,
		}, map[api.Endpoint]http.HandlerFunc{
			apicfg.Objectives.Private.RockCreate: priv.RockCreate,
		},
	)
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
