package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OpObjectiveCreate struct {
	Opid      OperationId
	Poid      ObjectiveId
	Pvid      VersionId
	Actor     UserId
	Content   string
	CreatedAt pgtype.Date
}

type OpObjectiveDelete struct {
	Opid      OperationId
	Oid       ObjectiveId
	Vid       VersionId
	Actor     UserId
	CreatedAt pgtype.Date
}

type OpObjectiveContentUpdate struct {
	Opid      OperationId
	Oid       ObjectiveId
	Vid       VersionId
	Actor     UserId
	Content   string
	CreatedAt pgtype.Date
}

type OpObjectiveAttachSubobjective struct {
	Opid      OperationId
	Actor     UserId
	SupOid    ObjectiveId
	SupVid    VersionId
	SubOid    ObjectiveId
	SubVid    VersionId
	CreatedAt pgtype.Date
}

type OpObjectiveUpdateCompletion struct {
	Opid      OperationId
	Oid       ObjectiveId
	Vid       VersionId
	Actor     UserId
	Completed string
	CreatedAt pgtype.Date
}

type Operation interface {
	op()
}

func (OpObjectiveCreate) op()
func (OpObjectiveDelete) op()
func (OpObjectiveContentUpdate) op()
func (OpObjectiveAttachSubobjective) op()
func (OpObjectiveUpdateCompletion) op()

func (db *Database) InsertOpObjectiveCreate(op OpObjectiveCreate) (OpObjectiveCreate, error) {
	const q = `
		INSERT INTO op_objective_create ("poid", "pvid", "actor", "content")
		VALUES ($1, $2, $3, $4)
		RETURNING ("opid", "poid", "pvid", "actor", "content", "created_at")
	`
	err := db.pool.QueryRow(context.Background(), q, op.Poid, op.Pvid, op.Actor, op.CreatedAt).Scan(
		&op.Opid, &op.Poid, &op.Pvid, &op.Actor, &op.Content, &op.CreatedAt,
	)
	if err != nil {
		return op, fmt.Errorf("query and scan: %w", err)
	}
	return op, nil
}

func (db *Database) InsertOpObjectiveDelete(op OpObjectiveDelete) (OpObjectiveDelete, error) {
	const q = `
		INSERT INTO op_objective_delete ("oid", "vid", "actor")
		VALUES ($1, $2, $3)
		RETURNING ("opid", "oid", "vid", "actor", "created_at")
	`
	err := db.pool.QueryRow(context.Background(), q, op.Oid, op.Vid, op.Actor).Scan(
		&op.Opid, &op.Oid, &op.Vid, &op.Actor, &op.CreatedAt,
	)
	if err != nil {
		return op, fmt.Errorf("query and scan: %w", err)
	}
	return op, nil
}

func (db *Database) InsertOpObjectiveContentUpdate(op OpObjectiveContentUpdate) (OpObjectiveContentUpdate, error) {
	const q = `
		INSERT INTO op_objective_content_update ("oid", "vid", "actor", "content")
		VALUES ($1, $2, $3, $4)
		RETURNING ("opid", "oid", "vid", "actor", "content", "created_at")
	`
	err := db.pool.QueryRow(context.Background(), q, op.Oid, op.Vid, op.Actor, op.Content).Scan(
		&op.Opid, &op.Oid, &op.Vid, &op.Actor, &op.Content, &op.CreatedAt,
	)
	if err != nil {
		return op, fmt.Errorf("query and scan: %w", err)
	}
	return op, nil
}

func (db *Database) InsertOpObjectiveAttachSubobjective(op OpObjectiveAttachSubobjective) (OpObjectiveAttachSubobjective, error) {
	const q = `
		INSERT INTO op_objective_attach_subobjective ("actor", "sup_oid", "sup_vid", "sub_oid", "sub_vid")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING ("opid", "actor", "sup_oid", "sup_vid", "sub_oid", "sub_vid", "created_at")
	`
	err := db.pool.QueryRow(context.Background(), q, op.Actor, op.SupOid, op.SupVid, op.SubOid, op.SubVid).Scan(
		&op.Opid, &op.Actor, &op.SupOid, &op.SupVid, &op.SubOid, &op.SubVid, &op.CreatedAt,
	)
	if err != nil {
		return op, fmt.Errorf("query and scan: %w", err)
	}
	return op, nil
}

func (db *Database) InsertOpObjectiveUpdateCompletion(op OpObjectiveUpdateCompletion) (OpObjectiveUpdateCompletion, error) {
	const q = `
		INSERT INTO op_objective_update_completion ("oid", "vid", "actor", "completed")
		VALUES ($1, $2, $3, $4)
		RETURNING ("opid", "oid", "vid", "actor", "completed", "created_at")
	`
	err := db.pool.QueryRow(context.Background(), q, op.Oid, op.Vid, op.Actor, op.Completed).Scan(
		&op.Opid, &op.Oid, &op.Vid, &op.Actor, &op.Completed, &op.CreatedAt,
	)
	if err != nil {
		return op, fmt.Errorf("query and scan: %w", err)
	}
	return op, nil
}
