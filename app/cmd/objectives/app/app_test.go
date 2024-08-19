package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"
	"logbook/models"
	"logbook/models/columns"
	"logbook/models/owners"
	"testing"
)

func TestApp(t *testing.T) {
	uid, err := columns.NewUuidV4[columns.UserId]()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, uid: %w", err))
	}

	srvcnf, err := service.ReadConfig("../local.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("reading service config: %w", err))
	}
	err = database.RunMigration(srvcnf)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}

	q, err := database.New(srvcnf.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, db connect: %w", err))
	}
	defer q.Close()

	a := New(q)
	ctx := context.Background()

	t.Run("rock", func(t *testing.T) {
		err = a.RockCreate(ctx, uid)
		if err != nil {
			t.Fatal(fmt.Errorf("act 1: %w", err))
		}
	})

	var rock models.Ovid
	t.Run("select bookmarks", func(t *testing.T) {
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
		rock.Vid, err = a.GetActiveVersion(ctx, rock.Oid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, GetActiveVersion: %w", err))
		}
	})

	t.Run("create first task", func(t *testing.T) {
		err := a.CreateSubtask(ctx, CreateSubtaskParams{
			Creator: uid,
			Parent:  rock,
			Content: "Hello world 1",
		})
		if err != nil {
			t.Fatal(fmt.Errorf("CreateSubtask: %w", err))
		}
	})

	var document []owners.DocumentItem
	t.Run("view build", func(t *testing.T) {
		rock.Vid, err = a.GetActiveVersion(ctx, rock.Oid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, GetActiveVersion: %w", err))
		}

		document, err = a.ViewBuilder(ctx, ViewBuilderParams{
			Viewer: uid,
			Root:   rock,
			Start:  0,
			Length: 2,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("ViewBuilder: %w", err))
		}

		if len(document) != 2 {
			t.Errorf("assert, document length, expected 2 got %d", len(document))
		}
	})

	t.Run("merged props", func(t *testing.T) {
		for _, e := range document {
			mps, err := a.GetMergedProps(ctx, models.Ovid{Oid: e.Oid, Vid: e.Vid})
			if err != nil {
				t.Fatal(fmt.Errorf("act, GetMergedProps: %w", err))
			}
			fmt.Println(e, mps)
		}
	})
}
