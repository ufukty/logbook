package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// objective or goal
type Objective struct {
	Oid      ObjectiveId
	Vid      VersionId
	Based    VersionId
	Content  string
	Creator  UserId
	Creation pgtype.Date
}

func (db *Database) InsertObjective(o Objective) (Objective, error) {
	q := `
		INSERT INTO "objective" ( "vid", "based", "content", "creator" ) 
		VALUES ( $1, $2, $3, $4) 
		RETURNING ( "oid", "creation" )`
	err := db.pool.QueryRow(context.Background(), q,
		&o.Vid, &o.Based, &o.Content, &o.Creator,
	).Scan(&o.Oid, &o.Creation)
	if err != nil {
		return o, fmt.Errorf("query and scan: %w", err)
	}
	return o, nil
}

func (db *Database) SelectObjective(ovid Ovid) (Objective, error) {
	q := `
		SELECT "oid", "vid", "based", "content", "creator", "creation"
		FROM "OBJECTIVE"
		WHERE "oid" = $1 AND "vid" == $2
		LIMIT 1`
	o := Objective{}
	err := db.pool.QueryRow(context.Background(), q, ovid.Oid, ovid.Vid).Scan(
		&o.Oid, &o.Vid, &o.Based, &o.Content, &o.Creator, &o.Creation,
	)
	if err != nil {
		return o, fmt.Errorf("query and scan: %w", err)
	}
	return o, nil
}

func (db *Database) SelectEffectiveVersionOfObjective(oid VersionId) (VersionId, error) {
	var vid VersionId
	q := `SELECT "vid" FROM "objective_effective_version" WHERE "oid" = $1 LIMIT 1`
	err := db.pool.QueryRow(context.Background(), q, oid).Scan(&vid)
	if err != nil {
		return "", fmt.Errorf("query & scan: %w", err)
	}
	return vid, nil
}
