package database

import (
	"logbook/internal/utilities/run"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.testing.env")
	os.Exit(m.Run())
}

func runMigration() {
	run.ExitAfterStderr("psql", "-U", os.Getenv("DBUSER"), "-d", os.Getenv("DBNAME_DEFAULT"),
		"-c", "DROP DATABASE IF EXISTS "+os.Getenv("DBNAME")+";",
		"-c", "CREATE DATABASE "+os.Getenv("DBNAME")+";")
	run.ExitAfterStderr("psql", "-U", os.Getenv("DBUSER"), "-d", os.Getenv("DBNAME"),
		"-f", "schema.sql")
}

func TestMigration(t *testing.T) {
	runMigration()
}
