package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/service"
	"logbook/internal/startup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Endpoints, error) {
	l, srvcnf, _, apicfg, err := startup.TestDependenciesWithServiceConfig("account", service.ReadConfig)
	if err != nil {
		return nil, fmt.Errorf("startup.TestDependenciesWithServiceConfig: %w", err)
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
	return New(a, l), nil
}
