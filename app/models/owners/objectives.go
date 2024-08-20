package owners

import (
	"logbook/cmd/objectives/queries"
	"logbook/models/columns"
	"time"
)

type Bookmark struct {
	Bid       columns.BookmarkId  `json:"bid"`
	Title     string              `json:"title"`
	Oid       columns.ObjectiveId `json:"oid"`
	IsRock    bool                `json:"is_rock"`
	CreatedAt time.Time           `json:"created_at"`
}

type ObjectiveType string

const (
	Goal = ObjectiveType("goal")
	Task = ObjectiveType("task")
)

type DocumentItem struct {
	Oid           columns.ObjectiveId `json:"oid"`
	Vid           columns.VersionId   `json:"vid"`
	Depth         int                 `json:"depth"`
	ObjectiveType ObjectiveType       `json:"objective_type"`
	Folded        bool                `json:"folded"`
}

type OperationHistoryItem struct {
	Version   columns.VersionId `json:"version"`
	Type      queries.OpType    `json:"type"`
	CreatedBy columns.UserId    `json:"created_by"`
	CreatedAt time.Time         `json:"created_at"`
}

type ObjectiveMergedProps struct {
	Content          string         `json:"content"`
	Completed        bool           `json:"completed"`
	SubtreeSize      int32          `json:"subtree_size"`
	SubtreeCompleted int32          `json:"subtree_completed"`
	Creator          columns.UserId `json:"creator"`
	Owner            columns.UserId `json:"owner"`
	CreatedAt        time.Time      `json:"created_at"`
}
