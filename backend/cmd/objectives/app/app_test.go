package app

import (
	"context"
	"encoding/json"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"
	"logbook/internal/startup"
	"logbook/internal/utils"
	"logbook/internal/utils/lines"
	"logbook/models"
	"logbook/models/columns"
	"logbook/models/owners"
	"math/rand/v2"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

type testfilenode struct {
	Content  columns.ObjectiveContent `json:"content"`
	Children []testfilenode           `json:"children"`
}

func testname(tc map[*testfilenode]*testfilenode) string {
	return fmt.Sprintf("registering %d objectives on %d parents", len(tc), len(utils.UniqueValues(tc)))
}

func TestAppManual(t *testing.T) {
	uid, err := columns.NewUuidV4[columns.UserId]()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, uid: %w", err))
	}

	l, srvcnf, _, err := startup.TestDependenciesWithServiceConfig("objectives", service.ReadConfig)
	if err != nil {
		t.Fatal(fmt.Errorf("startup.TestDependenciesWithServiceConfig: %w", err))
	}
	err = database.RunMigration(srvcnf)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, srvcnf.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("pgxpool.New: %w", err))
	}
	defer pool.Close()
	a := New(pool, l)

	t.Run("rock create", func(t *testing.T) {
		err = a.RockCreate(ctx, uid)
		if err != nil {
			t.Fatal(fmt.Errorf("act 1: %w", err))
		}
	})

	var rock models.Ovid
	t.Run("rock get", func(t *testing.T) {
		rock.Oid, err = a.RockGet(ctx, uid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, RockGet: %w", err))
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

const testFileNodes = 275 + 1 // number of nodes in the test file + rock

func TestDemoFile(t *testing.T) {
	a, uid, err := testdeps()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, testdeps: %w", err))
	}
	defer a.pool.Close()

	rock, err := loadDemo(context.Background(), a, uid, true)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, CreateDemoFileInDfsOrder: %w", err))
	}
	fmt.Println("rock:", rock)

	var firstGrandGrandChild models.Ovid // should be "Define Project Scope"
	t.Run("children", func(t *testing.T) {
		firstGrandGrandChild = rock
		for generation := range 3 {
			children, err := a.ListChildren(context.Background(), firstGrandGrandChild)
			if err != nil {
				t.Fatal(fmt.Errorf("ListChildren: %w", err))
			}
			if len(children) == 0 {
				t.Fatalf("no children from gen %d", generation)
			}
			fmt.Println("number of children:", len(children))
			firstGrandGrandChild = children[0]
		}
	})

	t.Run("checking props and bups of rock after loading document", func(t *testing.T) {
		rock.Vid, err = a.GetActiveVersion(context.Background(), rock.Oid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, GetActiveVersion: %w", err))
		}

		mp, err := a.GetMergedProps(context.Background(), rock)
		if err != nil {
			t.Fatal(fmt.Errorf("a.GetMergedProps: %w", err))
		}

		if mp.SubtreeSize == 0 {
			t.Fatal(fmt.Errorf("mp.SubtreeSize is expected to be bigger than 0"))
		}
	})

	t.Run("active path", func(t *testing.T) {
		ap, err := a.l2.ListActivePathToRock(context.Background(), a.oneshot, firstGrandGrandChild)
		if err != nil {
			t.Fatal(fmt.Errorf("listActivePathToRock: %w", err))
		}
		expected := 3 + 1
		if len(ap) != expected {
			t.Errorf("assert, expected %d got %d", expected, len(ap))
		}
	})

	var document []owners.DocumentItem
	t.Run("view build", func(t *testing.T) {
		ViewportLimit = 9999999

		document, err = a.ViewBuilder(context.Background(), ViewBuilderParams{
			Viewer: uid,
			Root:   rock,
			Start:  0,
			Length: ViewportLimit,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("ViewBuilder: %w", err))
		}

		if len(document) != testFileNodes {
			t.Errorf("assert, document length, expected %d got %d", testFileNodes, len(document))
		}
	})

	t.Run("merged props and print", func(t *testing.T) {
		o, err := os.Create(fmt.Sprintf("testresults/view_building_%s.txt", time.Now().Format("20060102_150405")))
		if err != nil {
			t.Fatal(fmt.Errorf("os.Create: %w", err))
		}
		defer o.Close()

		for _, e := range document {
			mps, err := a.GetMergedProps(context.Background(), models.Ovid{Oid: e.Oid, Vid: e.Vid})
			if err != nil {
				t.Fatal(fmt.Errorf("act, GetMergedProps: %w", err))
			}
			fmt.Fprintln(o, e, mps)
		}
	})

	t.Run("list children of root and delete them", func(t *testing.T) {
		rock.Vid, err = a.GetActiveVersion(context.Background(), rock.Oid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, GetActiveVersion: %w", err))
		}

		children, err := a.ListChildren(context.Background(), rock)
		if err != nil {
			t.Fatal(fmt.Errorf("ListChildren: %w", err))
		}

		for _, child := range children {
			fmt.Println("deleting subtask:", child)
			err := a.DeleteSubtask(context.Background(), DeleteSubtaskParams{
				Subject: child,
				Actor:   uid,
			})
			if err != nil {
				t.Fatal(fmt.Errorf("DeleteSubtask: %w", err))
			}

			rock.Vid, err = a.GetActiveVersion(context.Background(), rock.Oid)
			if err != nil {
				t.Fatal(fmt.Errorf("act, GetActiveVersion: %w", err))
			}

			props, err := a.GetMergedProps(context.Background(), rock)
			if err != nil {
				t.Fatal(fmt.Errorf("GetMergedProps: %w", err))
			}

			fmt.Println("subtree size of rock:", props.SubtreeSize)

			document, err = a.ViewBuilder(context.Background(), ViewBuilderParams{
				Viewer: uid,
				Root:   rock,
				Start:  0,
				Length: 2,
			})
			if err != nil {
				t.Fatal(fmt.Errorf("ViewBuilder: %w", err))
			}

			for _, e := range document {
				mps, err := a.GetMergedProps(context.Background(), models.Ovid{Oid: e.Oid, Vid: e.Vid})
				if err != nil {
					t.Fatal(fmt.Errorf("act, GetMergedProps: %w", err))
				}
				fmt.Println(e, mps)
			}
		}
	})

	t.Run("checking props and bups of rock after deleting document", func(t *testing.T) {
		rock.Vid, err = a.GetActiveVersion(context.Background(), rock.Oid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, GetActiveVersion: %w", err))
		}

		mp, err := a.GetMergedProps(context.Background(), rock)
		if err != nil {
			t.Fatal(fmt.Errorf("a.GetMergedProps: %w", err))
		}

		if mp.SubtreeSize != 0 {
			t.Fatal(fmt.Errorf("mp.SubtreeSize expected to be 0"))
		}
	})
}

func TestAppRandomOrderSubtaskCreationWithConcurrency(t *testing.T) {
	uid, err := columns.NewUuidV4[columns.UserId]()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, uid: %w", err))
	}

	l, srvcnf, _, err := startup.TestDependenciesWithServiceConfig("objectives", service.ReadConfig)
	if err != nil {
		t.Fatal(fmt.Errorf("startup.TestDependenciesWithServiceConfig: %w", err))
	}

	err = database.RunMigration(srvcnf)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, srvcnf.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("pgxpool.New: %w", err))
	}
	defer pool.Close()
	a := New(pool, l)

	t.Run("rock create", func(t *testing.T) {
		err = a.RockCreate(ctx, uid)
		if err != nil {
			t.Fatal(fmt.Errorf("act 1: %w", err))
		}
	})

	var rock models.Ovid
	t.Run("rock get", func(t *testing.T) {
		rock.Oid, err = a.RockGet(ctx, uid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, RockGet: %w", err))
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
						Creator: uid,
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
					store[child] = oid.Oid
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
