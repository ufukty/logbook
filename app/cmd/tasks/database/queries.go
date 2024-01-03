package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
)

func (db *Database) SelectPreviousVersion(oid ObjectiveId, vid VersionId) (VersionId, error) {
	q := `SELECT "prev" FROM "OPERATIONS" WHERE "vid" = $1 LIMIT 1`
	rows, err := db.pool.Query(context.Background(), q, vid)
	if err != nil {
		return "", fmt.Errorf("query: %w", err)
	}
	prev := new(pgtype.Text)
	err = rows.Scan(&prev)
	if err != nil {
		return "", fmt.Errorf("scan: %w", err)
	}
	if prev.Status != pgtype.Present {
		return "", fmt.Errorf("doesn't exist")
	}
	return VersionId(prev.String), nil
}

func (db *Database) InsertObjective(o Objective) (Objective, error) {
	q := `
		INSERT INTO "objective" ( "vid", "based", "type", "content", "creator" ) 
		VALUES ( $1, $2, $3, $4, $5 ) 
		RETURNING ( "oid", "creation" )`
	err := db.pool.QueryRow(context.Background(), q,
		&o.Vid, &o.Based, &o.Type, &o.Content, &o.Creator,
	).Scan(&o.Oid, &o.Creation)
	if err != nil {
		return o, fmt.Errorf("query and scan: %w", err)
	}
	return o, nil
}

func (db *Database) SelectObjective(oid ObjectiveId, vid VersionId) (Objective, error) {
	q := `
		SELECT "oid", "vid", "based", "type", "content", "creator", "creation"
		FROM "OBJECTIVE"
		WHERE "oid" = $1 AND "vid" == $2
		LIMIT 1`
	o := Objective{}
	err := db.pool.QueryRow(context.Background(), q, oid, vid).Scan(
		&o.Oid, &o.Vid, &o.Based, &o.Type, &o.Content, &o.Creator, &o.Creation,
	)
	if err != nil {
		return o, fmt.Errorf("query and scan: %w", err)
	}
	return o, nil
}

func (db *Database) SelectSubLinks(supoid string, supvid VersionId) ([]Link, error) {
	ls := []Link{}
	q := `
		SELECT "lid", "sup_oid", "sup_vid", "sub_oid", "sub_vid", "creation"
		FROM "objective_link" 
		WHERE "supoid" = $1 AND "supvid" = $2 
		LIMIT 50`
	rs, err := db.pool.Query(context.Background(), q, supoid, supvid)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	for rs.Next() {
		l := Link{}
		err := rs.Scan(
			&l.Lid, &l.SupOid, &l.SupVid, &l.SubOid, &l.SubVid, &l.Creation,
		)
		if err != nil {
			return ls, fmt.Errorf("scan: %w", err)
		}
		ls = append(ls, l)
	}
	return ls, nil
}

func (db *Database) SelectTheUpperLink(suboid ObjectiveId, subvid VersionId) (Link, error) {
	l := Link{}
	q := `
		SELECT "lid", "sup_oid", "sup_vid", "sub_oid", "sub_vid", "creation"
		FROM "objective_link" 
		WHERE "suboid" = $1 AND "subvid" = $2 
		LIMIT 1`
	err := db.pool.QueryRow(context.Background(), q, suboid, subvid).Scan(
		&l.Lid, &l.SupOid, &l.SupVid, &l.SubOid, &l.SubVid, &l.Creation,
	)
	if err != nil {
		return Link{}, fmt.Errorf("query and scan: %w", err)
	}
	return l, nil
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
