package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"
	"logbook/internal/logger"
	"logbook/models/columns"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func testdeps() (*App, columns.UserId, error) {
	l := logger.New("test")

	uid, err := columns.NewUuidV4[columns.UserId]()
	if err != nil {
		return nil, columns.ZeroUserId, fmt.Errorf("prep, uid: %w", err)
	}
	srvcnf, err := service.ReadConfig("../local.yml")
	if err != nil {
		return nil, columns.ZeroUserId, fmt.Errorf("reading service config: %w", err)
	}
	err = database.RunMigration(srvcnf)
	if err != nil {
		return nil, columns.ZeroUserId, fmt.Errorf("running migration: %w", err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, srvcnf.Database.Dsn)
	if err != nil {
		return nil, columns.ZeroUserId, fmt.Errorf("pgxpool.New: %w", err)
	}
	a := New(pool, l)
	return a, uid, nil
}

func TestGetObjectiveHistory(t *testing.T) {
	a, uid, err := testdeps()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, testdeps: %w", err))
	}
	defer a.pool.Close()

	rock, err := loadDemo(context.Background(), a, uid, true)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, CreateDemoFileInDfsOrder: %w", err))
	}
	children, err := a.ListChildren(context.Background(), rock)
	if err != nil {
		t.Fatal(fmt.Errorf("ListChildren: %w", err))
	}
	first := children[0]
	fmt.Println("choosing the child:", first.Oid, first.Vid)
	grandchildren, err := a.ListChildren(context.Background(), first)
	if err != nil {
		t.Fatal(fmt.Errorf("ListChildren/2: %w", err))
	}
	history, err := a.GetObjectiveHistory(context.Background(), GetObjectiveHistoryParams{first, false})
	if err != nil {
		t.Fatal(fmt.Errorf("GetObjectiveHistory: %w", err))
	}
	for _, item := range history {
		fmt.Println(item)
	}
	if len(history) != 1+len(grandchildren) {
		t.Errorf("len(history)=%d <> len(grandchildren)=%d", len(history), len(grandchildren))
	}
}
