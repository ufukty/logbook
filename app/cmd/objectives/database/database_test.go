package database

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/service"
	"logbook/models/columns"
	"testing"
)

func TestMigration(t *testing.T) {
	srvcfg, err := service.ReadConfig("../local.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("reading service config: %w", err))
	}
	err = RunMigration(srvcfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}
}

func TestLogic(t *testing.T) {
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

	uid, err := columns.NewUuidV4[columns.UserId]()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, generating new uid: %w", err))
	}

	ctx := context.Background()

	var rock Objective
	t.Run("create the rock", func(t *testing.T) {
		oid, err := columns.NewUuidV4[columns.ObjectiveId]()
		if err != nil {
			t.Fatal(fmt.Errorf("prep, generating new oid: %w", err))
		}
		vid, err := columns.NewUuidV4[columns.VersionId]()
		if err != nil {
			t.Fatal(fmt.Errorf("prep, generating new vid: %w", err))
		}

		register, err := q.InsertOperation(ctx, InsertOperationParams{
			Subjectoid: oid,
			Subjectvid: vid,
			Actor:      uid,
			OpType:     OpTypeUsrRegister,
			OpStatus:   OpStatusAccepted,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 1, insert operation for user registration (to create the rock with the opid): %w", err))
		}

		props, err := q.InsertProperties(ctx, InsertPropertiesParams{
			Content:   "",
			Completed: false,
			Creator:   uid,
			Owner:     uid,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 2, insert props: %w", err))
		}

		bup, err := q.InsertBottomUpProps(ctx, InsertBottomUpPropsParams{
			SubtreeSize:       0,
			CompletedSubitems: 0,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 3, insert bottom-up props: %w", err))
		}

		rock, err = q.InsertNewObjective(ctx, InsertNewObjectiveParams{
			CreatedBy: register.Opid,
			Pid:       props.Pid,
			Bupid:     bup.Bupid,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 3, insert the rock: %w", err))
		}
	})

	var obj Objective
	t.Run("generate the first task", func(t *testing.T) {
		op, err := q.InsertOperation(ctx, InsertOperationParams{
			Subjectoid: rock.Oid,
			Subjectvid: rock.Vid,
			Actor:      uid,
			OpType:     OpTypeObjCreateSubtask,
			OpStatus:   OpStatusReceived,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 1, insert row to 'operation' table: %w", err))
		}

		_, err = q.InsertOpObjCreateSubtask(ctx, InsertOpObjCreateSubtaskParams{
			Opid:    op.Opid,
			Content: "Hello world",
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 2, insert row to 'op_obj_create' table: %w", err))
		}

		props, err := q.InsertProperties(ctx, InsertPropertiesParams{
			Content:   "",
			Completed: false,
			Creator:   uid,
			Owner:     uid,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 3, insert props: %w", err))
		}

		bup, err := q.InsertBottomUpProps(ctx, InsertBottomUpPropsParams{
			SubtreeSize:       0,
			CompletedSubitems: 0,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 4, insert bottom-up props: %w", err))
		}

		obj, err = q.InsertNewObjective(ctx, InsertNewObjectiveParams{
			CreatedBy: op.Opid,
			Pid:       props.Pid,
			Bupid:     bup.Bupid,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 5, insert row to 'objective' table: %w", err))
		}
	})

	var rock2 Objective
	t.Run("link the task to rock", func(t *testing.T) {
		_, err := q.InsertLink(ctx, InsertLinkParams{
			SupOid: rock.Oid,
			SupVid: rock.Vid,
			SubOid: obj.Oid,
			SubVid: obj.Vid,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 1, insert row to 'link' table: %w", err))
		}

		l, err := q.SelectUpperLinks(ctx, SelectUpperLinksParams{
			SubOid: obj.Oid,
			SubVid: obj.Vid,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 2, select the link from objective to the rock: %w", err))
		}

		if len(l) != 1 {
			t.Fatal(fmt.Errorf("assert 1, len(l), expected 1, got %d", len(l)))
		}

		rock2, err = q.SelectObjective(ctx, SelectObjectiveParams{
			Oid: l[0].SupOid,
			Vid: l[0].SupVid,
		})
		if err != nil {
			t.Fatal(fmt.Errorf("act 3, select the rock through the link: %w", err))
		}
	})

	if rock != rock2 {
		t.Fatal(fmt.Errorf("assert 2,\n\trock1=%v\n\trock2=%v", rock, rock2))
	}

}
