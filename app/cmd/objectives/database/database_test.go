package database

import (
	"fmt"
	"logbook/cmd/objectives/service"
	"testing"
)

func TestMigration(t *testing.T) {
	srvcfg, err := service.ReadConfig("../testing.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("reading service config: %w", err))
	}
	err = RunMigration(srvcfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}
}
