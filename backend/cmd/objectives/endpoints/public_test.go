package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"
	sessions "logbook/cmd/sessions/client"
	"logbook/internal/startup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Public, error) {
	l, srvcnf, _, err := startup.TestDependenciesWithServiceConfig("objectives", service.ReadConfig)
	if err != nil {
		return nil, fmt.Errorf("startup.TestDependenciesWithServiceConfig: %w", err)
	}

	err = database.RunMigration(srvcnf)
	if err != nil {
		return nil, fmt.Errorf("running migration: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), srvcnf.Database.Dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	app := app.New(pool, l)
	ep := NewPublic(app, &sessions.Mock{}, l)
	return ep, nil
}
