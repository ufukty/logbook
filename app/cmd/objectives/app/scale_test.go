package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/cmd/objectives/service"
	"logbook/models"
	"logbook/models/columns"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestStorageScale(t *testing.T) {
	uid, err := columns.NewUuidV4[columns.UserId]()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, uid: %w", err))
	}
	srvcnf, err := service.ReadConfig("../local.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("reading service config: %w", err))
	}
	err = queries.RunMigration(srvcnf)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, srvcnf.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("pgxpool.New: %w", err))
	}
	defer pool.Close()
	a := New(pool)
	err = a.RockCreate(ctx, uid)
	if err != nil {
		t.Fatal(fmt.Errorf("act 1: %w", err))
	}
	var rock models.Ovid
	bs, err := a.ListBookmarks(ctx, ListBookmarksParams{Viewer: uid})
	if err != nil {
		t.Fatal(fmt.Errorf("listing bookmarks: %w", err))
	}
	found := false
	for _, b := range bs {
		if b.IsRock {
			rock = models.Ovid{Oid: b.Oid}
			found = true
		}
	}
	if !found {
		t.Fatal(fmt.Errorf("assert, expected rock to be found: %v", bs))
	}

	o, err := os.Create("scale_test.csv")
	if err != nil {
		t.Fatal(fmt.Errorf("os.Create: %w", err))
	}
	defer o.Close()

	store := []columns.ObjectiveId{rock.Oid}
	for i := range 1000 {
		parent := store[rand.IntN(len(store))]
		vid, err := a.GetActiveVersion(context.Background(), parent)
		if err != nil {
			t.Fatal(fmt.Errorf("GetActiveVersion(%d, %s): %w", i, parent, err))
		}
		oid, err := a.CreateSubtask(context.Background(), CreateSubtaskParams{
			Creator: uid,
			Parent: models.Ovid{
				Oid: parent,
				Vid: vid,
			},
			Content: strconv.Itoa(i),
		})
		if err != nil {
			t.Fatal(fmt.Errorf("CreateSubtask(%d, %s, %s): %w", i, parent, vid, err))
		}
		store = append(store, oid)

		r := pool.QueryRow(ctx, `SELECT count(*) FROM "objective"`)
		var count int
		if err := r.Scan(&count); err != nil {
			t.Fatal("querying the number of rows")
		}

		fmt.Fprintf(o, "%d;%d\n", i, count)
		time.Sleep(time.Millisecond * 50)
	}
}
