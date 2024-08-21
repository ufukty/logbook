package app

import (
	"context"
	"encoding/json"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/cmd/objectives/service"
	"logbook/internal/utilities/mapw"
	"logbook/internal/utilities/slicew/lines"
	"logbook/models"
	"logbook/models/columns"
	"logbook/models/owners"
	"math/rand/v2"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

type testfilenode struct {
	Content  string         `json:"content"`
	Children []testfilenode `json:"children"`
}

func testname(tc map[*testfilenode]*testfilenode) string {
	return fmt.Sprintf("registering %d objectives on %d parents", len(tc), len(mapw.UniqueValues(tc)))
}

func TestAppManual(t *testing.T) {
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
		_, err := a.CreateSubtask(ctx, CreateSubtaskParams{
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

func TestAppRandomOrderSubtaskCreation(t *testing.T) {
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

	testfile := []testfilenode{}
	t.Run("reading testdata file", func(t *testing.T) {
		f, err := os.Open("testdata/company.md.json")
		if err != nil {
			t.Fatal(fmt.Errorf("opening: %w", err))
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(&testfile)
		if err != nil {
			t.Fatal(fmt.Errorf("decoding: %w", err))
		}
	})

	ordered := map[*testfilenode]*testfilenode{} // []{node:parent}
	t.Run("randomizing ordering", func(t *testing.T) {
		visited_waiting := map[*testfilenode]*testfilenode{} // {node:parent}
		for _, tc := range testfile {
			visited_waiting[&tc] = nil // nil is for Rock
		}
		for len(visited_waiting) > 0 {
			rnd := rand.IntN(len(visited_waiting))
			child := maps.Keys(visited_waiting)[rnd]
			parent := visited_waiting[child]

			ordered[child] = parent
			for _, grandchild := range child.Children {
				visited_waiting[&grandchild] = child
			}

			delete(visited_waiting, child)
		}
	})

	store := map[*testfilenode]columns.ObjectiveId{
		nil: rock.Oid,
	}
	t.Run("creating subtasks", func(t *testing.T) {

		for child, parent := range ordered {
			parentOid, ok := store[parent]
			if !ok {
				t.Fatal(fmt.Errorf("registering %q on %q: test is shortcutting the hierarchy", child.Content, parent.Content))
			}

			vid, err := a.GetActiveVersion(context.Background(), parentOid)
			if err != nil {
				t.Fatal(fmt.Errorf("GetActiveVersion: %w", err))
			}
			oid, err := a.CreateSubtask(context.Background(), CreateSubtaskParams{
				Creator: columns.ZeroUserId,
				Parent: models.Ovid{
					Oid: parentOid,
					Vid: vid,
				},
				Content: child.Content,
			})
			if err != nil {
				t.Fatal(fmt.Errorf("CreateSubtask: %w", err))
			}
			store[child] = oid
		}
	})

	var document []owners.DocumentItem
	t.Run("view build", func(t *testing.T) {
		registered := len(store)

		rock.Vid, err = a.GetActiveVersion(ctx, rock.Oid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, GetActiveVersion: %w", err))
		}

		document, err = a.ViewBuilder(ctx, ViewBuilderParams{
			Viewer: uid,
			Root:   rock,
			Start:  0,
			Length: registered,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("ViewBuilder: %w", err))
		}

		if len(document) != registered {
			t.Errorf("assert, document length, expected %d got %d", registered, len(document))
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

func TestAppRandomOrderSubtaskCreationWithConcurrency(t *testing.T) {
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

	testfile := []testfilenode{}
	t.Run("reading testdata file", func(t *testing.T) {
		f, err := os.Open("testdata/company.md.json")
		if err != nil {
			t.Fatal(fmt.Errorf("opening: %w", err))
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(&testfile)
		if err != nil {
			t.Fatal(fmt.Errorf("decoding: %w", err))
		}
	})

	jobs := []map[*testfilenode]*testfilenode{{}} // []{node:parent}
	t.Run("preparing asynchronous calls", func(t *testing.T) {
		visited_waiting := map[*testfilenode]*testfilenode{} // {node:parent}
		for _, tc := range testfile {
			visited_waiting[&tc] = nil // nil is for Rock
		}
		for len(visited_waiting) > 0 {
			rnd := rand.IntN(len(visited_waiting))
			child := maps.Keys(visited_waiting)[rnd]
			parent := visited_waiting[child]

			_, same := jobs[len(jobs)-1][parent]        // can't send a parent and its one child at same job
			full := len(jobs[len(jobs)-1]) == len(jobs) // increasing number of asynchronous calls at each job
			if same || full {
				jobs = append(jobs, map[*testfilenode]*testfilenode{})
			}
			jobs[len(jobs)-1][child] = parent
			for _, grandchild := range child.Children {
				visited_waiting[&grandchild] = child
			}

			delete(visited_waiting, child)
		}
	})

	store := map[*testfilenode]columns.ObjectiveId{
		nil: rock.Oid,
	}
	for _, tc := range jobs {
		t.Run(testname(tc), func(t *testing.T) {
			var wg sync.WaitGroup
			var errs []string
			for child, parent := range tc {
				wg.Add(1)
				go func() {
					defer wg.Done()
					parentOid, ok := store[parent]
					if !ok {
						errs = append(errs, "test is shortcutting the hierarchy")
						return
					}
					vid, err := a.GetActiveVersion(context.Background(), parentOid)
					if err != nil {
						errs = append(errs, fmt.Sprintf("GetActiveVersion: %s", err.Error()))
						return
					}
					oid, err := a.CreateSubtask(context.Background(), CreateSubtaskParams{
						Creator: columns.ZeroUserId,
						Parent: models.Ovid{
							Oid: parentOid,
							Vid: vid,
						},
						Content: child.Content,
					})
					if err != nil {
						errs = append(errs, fmt.Sprintf("CreateSubtask: %s", err.Error()))
						return
					}
					store[child] = oid
				}()
			}
			wg.Wait()
			if len(errs) > 0 {
				t.Errorf("found %d errors:\n%s", len(errs), lines.Join(errs, "+ "))
			}
		})
	}

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
			Length: 250,
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
