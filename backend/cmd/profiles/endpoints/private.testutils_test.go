package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/profiles/app"
	"logbook/cmd/profiles/database"
	"logbook/cmd/profiles/service"
	"logbook/internal/startup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Private, error) {
	l, srvcnf, _, err := startup.TestDependenciesWithServiceConfig("profiles", service.ReadConfig)
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

	a := app.New(pool)
	return NewPrivate(a, l), nil
}
