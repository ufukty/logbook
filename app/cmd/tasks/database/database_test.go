package database

import (
	"logbook/internal/utilities/run"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.test.local.env")
	os.Exit(m.Run())
}

func runMigration() {
	run.ExitAfterStderr("psql", "-U", "ufuktan", "-d", "postgres", "-f", "migration.sql")
}

func TestMigration(t *testing.T) {
	runMigration()
}
