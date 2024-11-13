package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/sessions/app"
	"logbook/cmd/sessions/database"
	"logbook/cmd/sessions/service"
	"logbook/internal/startup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Public, error) {
	l, srvcnf, _, err := startup.TestDependenciesWithServiceConfig("sessions", service.ReadConfig)
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

	a := app.New(pool) // FIXME: mock objectives service?
	return NewPublic(a, l), nil
}
