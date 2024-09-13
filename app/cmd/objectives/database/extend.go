package database

import "fmt"

func lastsix[S ~string](id S) S {
	return id[max(0, len(id)-6):]
}

func (l Link) String() string {
	return fmt.Sprintln(
		"SupOid:", lastsix(l.SupOid),
		"SupVid:", lastsix(l.SupVid),
		"SubOid:", lastsix(l.SubOid),
		"SubVid:", lastsix(l.SubVid),
		"CreatedAtOriginal:", l.CreatedAtOriginal,
	)
}
