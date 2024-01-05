package app

import (
	"fmt"
	"logbook/cmd/tasks/database"
)

type ObjectiveVersionId struct {
	Oid database.ObjectiveId
	Vid database.VersionId
}

func (a *App) ListObjectiveAncestry(oid database.ObjectiveId, vid database.VersionId) ([]ObjectiveVersionId, error) {
	anc := []ObjectiveVersionId{ObjectiveVersionId{oid, vid}}
	c := anc[0]
	for limit := 0; true; limit++ {
		l, err := a.db.SelectTheUpperLink(c.Oid, c.Vid)
		if err != nil {
			return nil, fmt.Errorf("db.SelectTheUpperLink(%q, %q): %w", c.Oid, c.Vid, err)
		}
		c.Oid = l.SupOid
		c.Vid = l.SupVid
		if c.Oid == database.NullObjectId {
			break
		}
		anc = append(anc, c)
		if limit == 100 {
			return nil, fmt.Errorf("depth limit (100) is reached")
		}
	}
	return anc, nil
}
