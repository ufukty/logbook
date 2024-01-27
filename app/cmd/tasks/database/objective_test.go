package database

import (
	"fmt"
	"os"
	"testing"
)

func Test_Objectives(t *testing.T) {
	runMigration()

	db, err := New(os.Getenv("DSN"))
	if err != nil {
		t.Fatal(fmt.Errorf("prep, db connect: %w", err))
	}
	defer db.Close()

	o1, err := db.InsertObjective(Objective{
		Vid:     ZeroVersionId,
		Based:   ZeroVersionId,
		Content: "Hello world",
		Creator: ZeroUserId,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act 1: %w", err))
	}

	o2, err := db.SelectObjective(Ovid{o1.Oid, o1.Vid})
	if err != nil {
		t.Fatal(fmt.Errorf("act 2: %w", err))
	}

	if o1 != o2 {
		t.Fatal("assert, o1 != o2")
	}

	if o2.CreatedAt == ZeroDate {
		t.Fatal("assert 2, o2.CreatedAt is not populated")
	}

	if o2.Oid == ZeroObjectId {
		t.Fatal("assert 2, o2.CreatedAt is not populated")
	}
}
