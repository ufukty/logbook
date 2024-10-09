package app

import (
	"logbook/models"
	"logbook/models/columns"
	"slices"
	"testing"
)

func newOvid() models.Ovid {
	return models.Ovid{columns.NewUuidV4Unsafe[columns.ObjectiveId](), columns.NewUuidV4Unsafe[columns.VersionId]()}
}

func TestPopCommonActivePath(t *testing.T) {
	var (
		c = []models.Ovid{
			newOvid(),
			newOvid(),
			newOvid(),
			newOvid(),
		}
		l = []models.Ovid{
			newOvid(),
			newOvid(),
			newOvid(),
			newOvid(),
		}
		r = []models.Ovid{
			newOvid(),
			newOvid(),
			newOvid(),
			newOvid(),
		}
	)

	l2, r2, c2 := popCommonActivePath(slices.Concat(l, c), slices.Concat(r, c))

	if !slices.Equal(l, l2) {
		t.Errorf("assert l = l2\n\tl=%v\n\tl2=%v", l, l2)
	}
	if !slices.Equal(r, r2) {
		t.Errorf("assert r = r2\n\tr=%v\n\tr2=%v", r, r2)
	}
	if !slices.Equal(c, c2) {
		t.Errorf("assert c = c2\n\tc=%v\n\tc2=%v", c, c2)
	}
}
