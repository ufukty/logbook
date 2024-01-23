package database

import (
	"context"
	"fmt"
)

type VersioningConfig struct {
	Oid       ObjectiveId
	First     VersionId
	Effective VersionId
}

func (db *Database) SelectVersioningConfig(oid ObjectiveId) (VersioningConfig, error) {
	vc := VersioningConfig{}
	q := `SELECT "oid", "first", "effective" FROM "versioning_config" WHERE "oid" = $1 LIMIT 1`
	r, err := db.pool.Query(context.Background(), q, oid)
	if err != nil {
		return vc, fmt.Errorf("query: %w", err)
	}
	err = r.Scan(&vc)
	if err != nil {
		return vc, fmt.Errorf("scan: %w", err)
	}
	return vc, nil
}
