package database

import (
	"fmt"
	"logbook/cmd/tags/service"
	"os"
	"path/filepath"
	"testing"
)

func TestMigration(t *testing.T) {
	cfg, err := service.ReadConfig(filepath.Join(os.Getenv("WORKSPACE"), "backend/cmd/tags/local.yml"))
	if err != nil {
		fmt.Println(fmt.Errorf("reading service config: %w", err))
	}
	err = RunMigration(cfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}
}
