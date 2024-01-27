package database

import (
	"fmt"
	"os"
	"testing"
)

func Test_Links(t *testing.T) {
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

	o2, err := db.InsertObjective(Objective{
		Vid:     ZeroVersionId,
		Based:   ZeroVersionId,
		Content: "Quick brown fox",
		Creator: ZeroUserId,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act 2: %w", err))
	}

	li, err := db.InsertLink(Link{
		SupOid: o1.Oid,
		SupVid: o1.Vid,
		SubOid: o2.Oid,
		SubVid: o2.Vid,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act 3, adding link: %w", err))
	}

	if li.CreatedAt == ZeroDate {
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
