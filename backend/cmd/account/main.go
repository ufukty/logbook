package main

import (
	"context"
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/cmd/account/endpoints"
	"logbook/cmd/account/service"
	objectives "logbook/cmd/objectives/client"
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
	l, args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig("account", service.ReadConfig)
	if err != nil {
		return fmt.Errorf("reading configs: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
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
	sc := sidecar.New(reg, deplcfg, []models.Service{models.Objectives}, l)
	defer sc.Stop()

	objectives := objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives)))

	a := app.New(pool, apicfg, objectives)

	pub := endpoints.NewPublic(a, l)
	pri := endpoints.NewPrivate(a, l)

	agent := reception.NewAgent(deplcfg, l)
	err = agent.RegisterEndpoints(pub, pri)
	if err != nil {
		return fmt.Errorf("agent.RegisterEndpoints: %w", err)
	}

	// TODO: tls between services needs certs per host(name)
	err = router.StartServer(router.ServerParameters{
		Address:  args.PrivateNetworkIp,
		Port:     deplcfg.Ports.Accounts,
		Router:   deplcfg.Router,
		Service:  models.Account,
		ServeMux: agent.Mux(),
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
