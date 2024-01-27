package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Test_Objectives(t *testing.T) {
	runMigration()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		t.Fatal(fmt.Errorf("prep, db connect: %w", err))
	}
	defer pool.Close()

	q := New(pool)

	o1, err := q.InsertObjective(context.Background(), InsertObjectiveParams{
		Vid:     ZeroVersionId,
		Based:   ZeroVersionId,
		Content: "Hello world",
		Creator: ZeroUserId,
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
