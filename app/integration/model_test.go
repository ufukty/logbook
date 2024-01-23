package database

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	run("brew", "services", "start", "postgresql@15")
	run("psql", "-U", "ufuktan", "-d", "postgres", "-f", "../cmd/tasks/database/migration.sql")

	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	os, err := load()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, load: %w", err))
	}

	if len(os) == 0 {
		t.Fatal(fmt.Errorf("prep, assert: test file has no objective instance to create"))
	}

	// create the Rock and get its id
	rockId, err := createTheRock()
	if err != nil {
		t.Fatal(fmt.Errorf("act, creating the rock: %w", err))
	}

	if err := createOnRock(rockId, os); err != nil {
		t.Fatal(fmt.Errorf("act, creating the rock: %w", err))
	}

}
