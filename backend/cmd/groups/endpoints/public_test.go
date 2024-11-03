package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/groups/app"
	"logbook/cmd/groups/database"
	"logbook/cmd/groups/service"
	"logbook/internal/startup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Public, error) {
	l, srvcnf, _, _, err := startup.TestDependenciesWithServiceConfig("groups", service.ReadConfig)
	if err != nil {
		return nil, fmt.Errorf("startup: %w", err)
	}
	err = database.Migrate(srvcnf)
	if err != nil {
		return nil, fmt.Errorf("running migration: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), srvcnf.Database.Dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	app := app.New(pool)
	ep := NewPublic(app, l)
	return ep, nil
}
