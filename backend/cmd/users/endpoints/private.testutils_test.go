package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/users/app"
	"logbook/cmd/users/database"
	"logbook/cmd/users/service"
	"logbook/internal/startup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Private, error) {
	l, srvcnf, _, err := startup.TestDependenciesWithServiceConfig("users", service.ReadConfig)
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
