package integration

import (
	"fmt"
	"logbook/integration/data"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	run("brew", "services", "start", "postgresql@15")
	run("psql", "-U", "ufuktan", "-d", "postgres", "-f", "../cmd/tasks/database/migration.sql")

	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	objs, err := data.LoadTestData()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, load: %w", err))
	}

	if len(objs) == 0 {
		t.Fatal(fmt.Errorf("prep, assert: test file has no objective instance to create"))
	}

	uctl, err := NewUserClient("../testing")
	if err != nil {
		t.Fatal(fmt.Errorf("prep, user client: %w", err))
	}

	if err = uctl.Register(); err != nil {
		t.Fatal(fmt.Errorf("act, registering: %w", err))
	}

	// create the Rock and get its id
	rockId, err := uctl.createTheRock()
	if err != nil {
		t.Fatal(fmt.Errorf("act, creating the rock: %w", err))
	}

	if err := createOnRock(rockId, objs); err != nil {
		t.Fatal(fmt.Errorf("act, creating the rock: %w", err))
	}

}
