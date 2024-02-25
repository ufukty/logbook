package objectives

import (
	"fmt"
	"logbook/cmd/tasks/app"
	"logbook/cmd/tasks/endpoints"
	"logbook/internal/utilities/run"
	"logbook/internal/web/api"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.test.local.env")
	os.Exit(m.Run())
}

func runMigration() {
	run.ExitAfterStderr("psql",
		"-U", os.Getenv("DBUSER"),
		"-d", os.Getenv("DBNAME_DEFAULT"),
		"-c", "DROP DATABASE IF EXISTS "+os.Getenv("DBNAME")+";",
		"-c", "CREATE DATABASE "+os.Getenv("DBNAME")+";")
	run.ExitAfterStderr("psql",
		"-U", os.Getenv("DBUSER"),
		"-d", os.Getenv("DBNAME"),
		"-f", "../database/schema.sql")
}

func TestMigration(t *testing.T) {
	runMigration()
}

func getDependencies() (*Client, error) {
	cfg, err := api.ReadConfig("../../../api.yml")
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	ctl := NewClient(cfg)

	// q, err := database.New(os.Getenv("DSN"))
	// if err != nil {
	// 	return nil, fmt.Errorf("prep, db connect: %w", err)
	// }
	// defer q.Close()

	// app := app.New(q)
	// ep := endpoints.NewManager(app)

	return ctl, nil
}

func TestDependencies(t *testing.T) {
	_, err := getDependencies()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartServer(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, server: %w", err))
	}
	defer s.Close()

	c, err := getDependencies()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, dependencies: %w", err))
	}

	bs, err := c.CreateObjective(&endpoints.CreateTaskRequest{
		Parent: app.Ovid{
			Oid: "00000000-0000-0000-0000-000000000001",
			Vid: "00000000-0000-0000-0000-000000000002",
		},
		Content: "Lorem ipsum dolor sit amet.",
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act, sending request: %w", err))
	}

	if len(bs.Update) == 0 {
		t.Fatalf("response.Update length is 0")
	}

}

func TestCreateObjectives(t *testing.T) {
	ctl, err := getDependencies()
	if err != nil {
		t.Fatal(fmt.Errorf("prep: %w", err))
	}

	bs, err := ctl.CreateObjective(&endpoints.CreateTaskRequest{
		Parent: app.Ovid{
			Oid: "",
			Vid: "",
		},
		Content: "",
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}

	if len(bs.Update) == 0 {
		t.Fatal("assert 1")
	}
}
