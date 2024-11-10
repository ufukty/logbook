package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"
	"logbook/internal/startup"
	"logbook/models"
	"logbook/models/columns"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func expIntN(n int) int {
	return int(float64(n) * rand.ExpFloat64() / math.MaxFloat64)
}

func TestStorageScale(t *testing.T) {
	uid, err := columns.NewUuidV4[columns.UserId]()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, uid: %w", err))
	}
	l, srvcnf, _, _, err := startup.TestDependenciesWithServiceConfig("objectives", service.ReadConfig)
	if err != nil {
		t.Fatal(fmt.Errorf("startup.TestDependenciesWithServiceConfig: %w", err))
	}
	err = database.RunMigration(srvcnf)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, running migration: %w", err))
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, srvcnf.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, pgxpool.New: %w", err))
	}
	defer pool.Close()
	a := New(pool, l)

	t.Run("rock create", func(t *testing.T) {
		err = a.RockCreate(ctx, uid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, RockCreate: %w", err))
		}
	})

	var rock models.Ovid
	t.Run("rock get", func(t *testing.T) {
		rock.Oid, err = a.RockGet(ctx, uid)
		if err != nil {
			t.Fatal(fmt.Errorf("act, RockGet: %w", err))
		}
	})

	o, err := os.Create(fmt.Sprintf("testresults/scale_test_%s.csv", time.Now().Format("20060102_150405")))
	if err != nil {
		t.Fatal(fmt.Errorf("os.Create: %w", err))
	}
	defer o.Close()

	store := []columns.ObjectiveId{rock.Oid}
	for i := range 1000 {
		parent := store[expIntN(len(store))]
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
			Content: columns.ObjectiveContent(strconv.Itoa(i)),
		})
		if err != nil {
			t.Fatal(fmt.Errorf("CreateSubtask(%d, %s, %s): %w", i, parent, vid, err))
		}
		store = append(store, oid.Oid)

		r := pool.QueryRow(ctx, `SELECT count(*) FROM "objective"`)
		var count int
		if err := r.Scan(&count); err != nil {
			t.Fatal("querying the number of rows")
		}

		fmt.Fprintf(o, "%d;%d\n", i, count)
	}
}
