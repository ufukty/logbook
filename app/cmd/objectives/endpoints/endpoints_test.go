package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getTestDependencies() (*Endpoints, error) {
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
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	app := app.New(pool)
	ep := New(app)
	return ep, nil
}
