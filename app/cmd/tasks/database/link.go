package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
)

type Link struct {
	Lid LinkId

	SupOid ObjectiveId
	SupVid VersionId
	SubOid ObjectiveId
	SubVid VersionId

	Creation pgtype.Date
}

func (db *Database) SelectSubLinks(supoid ObjectiveId, supvid VersionId) ([]Link, error) {
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

func (db *Database) SelectTheUpperLink(sub Ovid) (Link, error) {
	l := Link{}
	q := `
		SELECT "lid", "sup_oid", "sup_vid", "sub_oid", "sub_vid", "creation"
		FROM "objective_link" 
		WHERE "suboid" = $1 AND "subvid" = $2 
		LIMIT 1`
	err := db.pool.QueryRow(context.Background(), q, sub.Oid, sub.Vid).Scan(
		&l.Lid, &l.SupOid, &l.SupVid, &l.SubOid, &l.SubVid, &l.Creation,
	)
	if err != nil {
		return Link{}, fmt.Errorf("query and scan: %w", err)
	}
	return l, nil
}

func (db *Database) InsertLink(l Link) (Link, error) {
	q := `
		INSERT INTO "objective_link" ( "sup_oid", "sup_vid", "sub_oid", "sub_vid" ) 
		VALUES ( $1, $2, $3, $4 ) 
		RETURNING ( "lid", "creation" )`
	err := db.pool.QueryRow(context.Background(), q,
		l.SubOid, l.SubVid, l.SupOid, l.SupVid,
	).Scan(&l.Lid, &l.Creation)
	if err != nil {
		return l, fmt.Errorf("query and scan: %w", err)
	}
	return l, nil
}
