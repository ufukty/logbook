package owners

import (
	"logbook/models/columns"
	"time"
)

type Bookmark struct {
	Bid       columns.BookmarkId  `json:"bid"`
	Title     string              `json:"title"`
	Oid       columns.ObjectiveId `json:"oid"`
	Vid       columns.VersionId   `json:"vid"`
	IsRock    bool                `json:"is_rock"`
	CreatedAt time.Time           `json:"created_at"`
}

type ObjectiveType string

const (
	Goal = ObjectiveType("goal")
	Task = ObjectiveType("task")
)

type ObjectiveView struct {
	Oid           columns.ObjectiveId `json:"oid"`
	Vid           columns.VersionId   `json:"vid"`
	Depth         int                 `json:"depth"`
	ObjectiveType ObjectiveType       `json:"objective_type"`
	Folded        bool                `json:"folded"`
}
