package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
)

func (db *Database) GetPreviousVersion(vid VersionId) (*VersionId, error) {
	q := `SELECT "prev" FROM "OPERATIONS" WHERE "vid" = $1 LIMIT 1`
	rows, err := db.pool.Query(context.Background(), q, vid)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	prev := new(pgtype.Text)
	err = rows.Scan(&prev)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	if prev.Status != pgtype.Present {
		return nil, fmt.Errorf("doesn't exist")
	}
	return (*VersionId)(&prev.String), nil
}

// Returns a list of updated objectives in addition to
// the created task in first item.
func (db *Database) CreateObjective(task *Objective) ([]Objective, error) {
	query := `
		SELECT
			"task_id",
			"document_id",
			"parent_id",
			"content",
			"degree",
			"depth",
			"created_at",
			"completed_at",
			"ready_to_pick_up"
		FROM
			create_task($1, $2, $3)`
	rows, err := db.pool.Query(
		context.Background(),
		query,
		&task.DocumentId,
		&task.Content,
		&task.ParentId,
	)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
	}
	tasks := []Objective{}
	for rows.Next() {
		task := Objective{}
		err := rows.Scan(
			&task.ObjectiveId,
			&task.DocumentId,
			&task.ParentId,
			&task.Content,
			&task.Degree,
			&task.Depth,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.ReadyToPickUp,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Database) GetObjective(oid ObjectiveId) (*Objective, error) {
	q := `
	SELECT 
		"oid", 
		"poid", 
		"creator", 
		"text", 
		"created_at", 
		"completed_at", 
		"archived_at"
	FROM 
		"OBJECTIVE" 
	WHERE 
		"oid" = $1
	LIMIT 1`
	r, err := db.pool.Query(context.Background(), q, oid)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	o := &Objective{}
	if err = r.Scan(&o.Oid, &o.ParentId, &o.Vid, &o.Creator, &o.Text, &o.CreatedAt, &o.CompletedAt, &o.ArchivedAt); err != nil {
		return nil, fmt.Errorf("scanning: %w", err)
	}
	return o, nil
}

func (db *Database) GetSubObjectives(oid string) ([]Objective, error) {
	tasks := []Objective{}
	query := `
	SELECT
			"task_id",
			"document_id",
			"parent_id",
			"content",
			"degree",
			"depth",
			"created_at",
			"completed_at",
			"ready_to_pick_up"
		FROM
			"TASK"
		WHERE
			"parent_id"=$1`
	rows, err := db.pool.Query(context.Background(), query, oid)
	if err != nil {
		return tasks, fmt.Errorf("running the query: %w", err)
	}
	for rows.Next() {
		task := Objective{}
		err = rows.Scan(&task.ObjectiveId, &task.DocumentId, &task.ParentId, &task.Content, &task.Degree, &task.Depth, &task.CreatedAt, &task.CompletedAt, &task.ReadyToPickUp)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	if err != nil {
		return nil, fmt.Errorf("scanning query results: %w", err)
	}
	return tasks, nil
}

// Returns a list of updated objectives in addition to
// the reattached task in first item.
func (db *Database) ReattachObjective(taskId string, newParentId string) ([]Objective, error) {
	query := `SELECT reattach_task($1, $2)`
	rows, err := db.pool.Query(context.Background(), query, taskId, newParentId)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
	}
	tasks := []Objective{}
	for rows.Next() {
		task := Objective{}
		err := rows.Scan(
			&task.ObjectiveId,
			&task.DocumentId,
			&task.ParentId,
			&task.Content,
			&task.Degree,
			&task.Depth,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.ReadyToPickUp,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Database) GetDocumentOverviewWithDocumentId(documentId string) ([]Objective, error) {
	tasks := []Objective{}
	query := `
		SELECT
			"task_id",
			"document_id",
			"parent_id",
			"content",
			"degree",
			"depth",
			"created_at",
			COALESCE("completed_at", '0001-01-01'),
			"ready_to_pick_up"
		FROM document_overview($1)`
	rows, err := db.pool.Query(context.Background(), query, documentId)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
	}
	for rows.Next() {
		task := Objective{}
		rows.Scan(
			&task.ObjectiveId,
			&task.DocumentId,
			&task.ParentId,
			&task.Content,
			&task.Degree,
			&task.Depth,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.ReadyToPickUp,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Database) GetChronologicalViewItems(documentId string, limit int, offset int) ([]Objective, error) {
	tasks := []Objective{}
	query := `
		SELECT
			"task_id",
			"document_id",
			"parent_id",
			"content",
			"degree",
			"depth",
			"created_at",
			"completed_at",
			"ready_to_pick_up"
		FROM "TASK"
		WHERE "document_id" = $1
		ORDER BY "created_at" ASC
		LIMIT $2 OFFSET $3
		`
	rows, err := db.pool.Query(context.Background(), query, documentId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
	}
	for rows.Next() {
		task := Objective{}
		err := rows.Scan(
			&(task.ObjectiveId),
			&(task.DocumentId),
			&(task.ParentId),
			&(task.Content),
			&(task.Degree),
			&(task.Depth),
			&(task.CreatedAt),
			&(task.CompletedAt),
			&(task.ReadyToPickUp),
		)
		if err != nil {
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Database) SuperLinksVersioned(oid ObjectiveId, vid VersionId) []Link {
	panic("not implemented")
}

func (db *Database) ListSuperObjectivesVersioned(oid ObjectiveId, vid VersionId) ([]Objective, error) {
	q := `SELECT * FROM "LINK" WHERE suboid = $1 AND subvid = $2`
	rows, err := db.pool.Query(context.Background(), q, oid, vid)
	if err != nil {
		return nil, fmt.Errorf("querying: %w", err)
	}
	os := []Objective{}
	for rows.Next() {
		o := Objective{}
		err := rows.Scan(&o.Oid, &o.ParentId, &o.Vid, &o.Creator, &o.Text, &o.CreatedAt, &o.CompletedAt, &o.ArchivedAt)
		if err != nil {
			return nil, fmt.Errorf("scannnig a row: %w", err)
		}
		os = append(os, o)
	}
	return os, nil
}
