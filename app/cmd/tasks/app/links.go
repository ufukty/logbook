package app

import (
	"fmt"
	"logbook/cmd/tasks/database"
)

func (a *App) ListObjectiveAncestry(ovid database.Ovid) ([]database.Ovid, error) {
	anc := []database.Ovid{}
	c := ovid
	for limit := 0; true; limit++ {
		l, err := a.queries.SelectTheUpperLink(c)
		if err != nil {
			return nil, fmt.Errorf("db.SelectTheUpperLink(%q, %q): %w", c.Oid, c.Vid, err)
		}
		c.Oid = l.SupOid
		c.Vid = l.SupVid
		if c.Oid == database.ZeroObjectId {
			break
		}
		anc = append(anc, c)
		if limit == 100 {
			return nil, fmt.Errorf("depth limit (100) is reached")
		}
	}
	return anc, nil
}
