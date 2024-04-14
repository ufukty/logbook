package database

import (
	"fmt"
	"logbook/cmd/account/service"
	"testing"
)

func TestMigration(t *testing.T) {
	cfg, err := service.ReadConfig("../testing.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("reading service config: %w", err))
	}
	err = RunMigration(cfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}
}
