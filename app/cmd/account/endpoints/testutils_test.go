package endpoints

import (
	"fmt"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/service"
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

	q, err := database.New(srvcnf.Database.Dsn)
	if err != nil {
		return nil, fmt.Errorf("connecting database: %w", err)
	}
	app := app.New(q)
	return New(app), nil
}
