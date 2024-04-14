package database

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
	"logbook/cmd/objectives/service"
	"testing"
)

func Test_Links(t *testing.T) {
	srvcnf, err := service.ReadConfig("../testing.yml")
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

	o2, err := q.InsertObjective(context.Background(), InsertObjectiveParams{
		Vid:     ZeroVersionId,
		Based:   ZeroVersionId,
		Content: "Quick brown fox",
		Creator: database.ZeroUserId,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act 2: %w", err))
	}

	li, err := q.InsertLink(context.Background(), InsertLinkParams{
		SupOid:  o1.Oid,
		SupVid:  o1.Vid,
		SubOid:  o2.Oid,
		SubVid:  o2.Vid,
		Creator: database.ZeroUserId,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act 3, adding link: %w", err))
	}

	if li.CreatedAt == ZeroTimestamp {
		t.Fatal(fmt.Println("assert 1, created_at is not populated"))
	}

	if li.Lid == ZeroLinkId {
		t.Fatal(fmt.Println("assert 2, lid is not populated"))
	}

	if li.SupOid != o1.Oid {
		t.Fatal(fmt.Println("assert 3, unexpected sup_oid"))
	}

	if li.SupVid != o1.Vid {
		t.Fatal(fmt.Println("assert 4, unexpected sup_vid"))
	}

	if li.SubOid != o2.Oid {
		t.Fatal(fmt.Println("assert 5, unexpected sub_oid"))
	}

	if li.SubVid != o2.Vid {
		t.Fatal(fmt.Println("assert 6, unexpected sub_vid"))
	}
}
