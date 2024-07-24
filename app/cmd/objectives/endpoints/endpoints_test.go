package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"
	"logbook/models"
)

type mockInstanceSource []models.Instance

func (m *mockInstanceSource) Instances() ([]models.Instance, error) {
	return *m, nil
}

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
	app := app.New(q, &mockInstanceSource{}) // FIXME:
	return New(app), nil
}
