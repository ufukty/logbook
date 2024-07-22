package endpoints

import (
	"fmt"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/service"
	"logbook/config/api"
	"logbook/models"
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

	q, err := database.New(srvcnf.Database.Dsn)
	if err != nil {
		return nil, fmt.Errorf("connecting database: %w", err)
	}
	a := app.New(q, apicfg, &MockInstanceSource{models.Instance{}}) // FIXME:
	return New(a), nil
}
