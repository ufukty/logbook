package database

import (
	"fmt"
	"logbook/cmd/groups/service"
	"testing"
)

func TestMigration(t *testing.T) {
	srvcfg, err := service.ReadConfig("../local.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("reading service config: %w", err))
	}
	err = Migrate(srvcfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}
}
