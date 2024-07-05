package database

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
	"logbook/cmd/objectives/service"
	"testing"
)

func Test_Objectives(t *testing.T) {
	srvcnf, err := service.ReadConfig("../local.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("reading service config: %w", err))
	}
	err = RunMigration(srvcnf)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}

	q, err := New(srvcnf.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, db connect: %w", err))
	}
	defer q.Close()

	o1, err := q.InsertObjective(context.Background(), InsertObjectiveParams{
		Vid:     ZeroVersionId,
		Based:   ZeroVersionId,
		Content: "Hello world",
		Creator: database.ZeroUserId,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act 1: %w", err))
	}

	o2, err := q.SelectObjective(context.Background(), SelectObjectiveParams{
		Oid: o1.Oid,
		Vid: o1.Vid,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act 2: %w", err))
	}

	if o1 != o2 {
		t.Fatal("assert, o1 != o2")
	}

	if o2.CreatedAt == ZeroTimestamp {
		t.Fatal("assert 2, o2.CreatedAt is not populated")
	}

	if o2.Oid == ZeroObjectId {
		t.Fatal("assert 2, o2.CreatedAt is not populated")
	}
}
