package database

import (
	"fmt"
	"logbook/cmd/objectives/service"
	"logbook/internal/utilities/run"
	"testing"
)

func runMigration(cfg service.Config) error {
	output := run.ExitAfterStderr("psql",
		"-U", cfg.Database.User,
		"-d", cfg.Database.Default,
		"-c", "DROP DATABASE IF EXISTS "+cfg.Database.Name+";",
		"-c", "CREATE DATABASE "+cfg.Database.Name+";",
	)
	if output != "" {
		return fmt.Errorf("dropping and recreating the application database: %s", output)
	}
	output = run.ExitAfterStderr("psql",
		"-U", cfg.Database.User,
		"-d", cfg.Database.Name,
		"-f", "schema.sql",
	)
	if output != "" {
		return fmt.Errorf("building the application database: %s", output)
	}
	return nil
}

func TestMigration(t *testing.T) {
	cfg, err := service.ReadConfig("../testing.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("reading service config: %w", err))
	}
	err = runMigration(cfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}
}
