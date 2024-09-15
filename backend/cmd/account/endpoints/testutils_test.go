package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/service"
	"logbook/config/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Endpoints, error) {
	apicfg, err := api.ReadConfig("../../../api.yml")
	if err != nil {
		return nil, fmt.Errorf("reading api config: %w", err)
	}
	srvcnf, err := service.ReadConfig("../local.yml")
	if err != nil {
		return nil, fmt.Errorf("reading service config: %w", err)
	}
	err = database.RunMigration(srvcnf)
	if err != nil {
		return nil, fmt.Errorf("running migration: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), srvcnf.Database.Dsn)
	if err != nil {
		return nil, fmt.Errorf("connecting database: %w", err)
	}
	defer pool.Close()

	a := app.New(pool, apicfg, nil) // FIXME: mock objectives service?
	return New(a), nil
}
