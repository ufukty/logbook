package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/database"
	"logbook/internal/utilities/run"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.test.local.env")
	os.Exit(m.Run())
}

func getTestDependencies() (*Endpoints, error) {
	run.ExitAfterStderr("psql", "-U", os.Getenv("DBUSER"), "-d", os.Getenv("DBNAME_DEFAULT"),
		"-c", "DROP DATABASE IF EXISTS "+os.Getenv("DBNAME")+";",
		"-c", "CREATE DATABASE "+os.Getenv("DBNAME")+";")
	run.ExitAfterStderr("psql", "-U", os.Getenv("DBUSER"), "-d", os.Getenv("DBNAME"),
		"-f", "../database/schema.sql")

	q, err := database.New(os.Getenv("DSN"))
	if err != nil {
		return nil, fmt.Errorf("connecting database: %w", err)
	}
	app := app.New(q)
	return NewManager(app), nil
}

func TestGetTestDependencies(t *testing.T) {
	_, err := getTestDependencies()
	if err != nil {
		t.Fatal(err)
	}

}
