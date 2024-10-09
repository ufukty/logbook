package main

import (
	"context"
	"fmt"
	"log"
	"logbook/cmd/account/api/private"
	"logbook/cmd/account/api/public"
	"logbook/cmd/account/service"
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
	l := logger.New("account")

	args, srvcfg, deplcfg, apicfg, err := startup.ServiceWithCustomConfig(service.ReadConfig, l)
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
	sc := sidecar.New(registry.NewClient(balancer.New(internalsd), apicfg, true), deplcfg, []models.Service{
		models.Objectives,
	}, l)
	defer sc.Stop()

	pub := public.New(apicfg, deplcfg, pool, sc, l)
	pri := private.New(apicfg, deplcfg, pool, l)

	// TODO: tls between services needs certs per host(name)
	router.StartServer(router.ServerParameters{
		Address: args.PrivateNetworkIp,
		Port:    deplcfg.Ports.Accounts,
		Router:  deplcfg.Router,
		Service: models.Account,
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
		log.Fatalln(err)
	}
}
