package database

import (
	"fmt"
	"logbook/internal/utilities/run"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	run.ExitAfterStderr("psql", "-U", "ufuktan", "-d", "postgres", "-f", "../database/migration.sql")
	os.Exit(m.Run())
}

func ExampleTestMain() {
	fmt.Println("Hello world.") // Output: Hello world.
}
