package database

import (
	"context"
	"fmt"
)

type Version struct {
	Vid   VersionId
	Based VersionId
}

func (db *Database) SelectVersion(vid VersionId) (Version, error) {
	v := Version{}
	q := `SELECT "vid", "based" FROM "version" WHERE "vid" = $1 LIMIT 1`
	rows, err := db.pool.Query(context.Background(), q, vid)
	if err != nil {
		return v, fmt.Errorf("query: %w", err)
	}
	err = rows.Scan(&v.Vid, &v.Based)
	if err != nil {
		return v, fmt.Errorf("scan: %w", err)
	}
	return v, nil
}

func (db *Database) InsertVersion(based VersionId) (Version, error) {
	v := Version{Based: based}
	q := `
		INSERT INTO "version" ( "based" ) 
		VALUES ( $1 ) 
		RETURNING ( "vid" )`
	err := db.pool.QueryRow(context.Background(), q, based).Scan(&v.Vid)
	if err != nil {
		return v, fmt.Errorf("query and scan: %w", err)
	}
	return v, nil
}

// func (db *Database) SelectPreviousVersion(oid ObjectiveId, vid VersionId) (VersionId, error) {
// 	q := `SELECT "prev" FROM "OPERATIONS" WHERE "vid" = $1 LIMIT 1`
// 	rows, err := db.pool.Query(context.Background(), q, vid)
// 	if err != nil {
// 		return "", fmt.Errorf("query: %w", err)
// 	}
// 	prev := new(pgtype.Text)
// 	err = rows.Scan(&prev)
// 	if err != nil {
// 		return "", fmt.Errorf("scan: %w", err)
// 	}
// 	if prev.Status != pgtype.Present {
// 		return "", fmt.Errorf("doesn't exist")
// 	}
// 	return VersionId(prev.String), nil
// }
