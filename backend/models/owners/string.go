package owners

import (
	"fmt"
	"logbook/internal/utilities/strw"
)

func ternary[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

func lastsix[S ~string](id S) S {
	return id[max(0, len(id)-6):]
}

func (ow DocumentItem) String() string {
	return fmt.Sprintf("%s(type:%s) (oid:%s) (vid:%s)%s",
		strw.Fill("  ", ow.Depth),
		ow.ObjectiveType,
		lastsix(ow.Oid),
		lastsix(ow.Vid),
		ternary(ow.Folded, " (fold)", ""),
	)
}

func (omps ObjectiveMergedProps) String() string {
	return fmt.Sprintf("(content:%s) (%s) (subtree:%d/%d) (owner:%s) (creator:%s) (%s)",
		omps.Content,
		ternary(omps.Completed, "completed", "todo"),
		omps.SubtreeCompleted,
		omps.SubtreeSize,
		lastsix(omps.Owner),
		lastsix(omps.Creator),
		omps.CreatedAt,
	)
}

func (hist OperationHistoryItem) String() string {
	return fmt.Sprintf("(Version:%s) (Type:%s) (CreatedBy:%s) (CreatedAt:%s)",
		lastsix(hist.Version),
		hist.Type,
		lastsix(hist.CreatedBy),
		hist.CreatedAt,
	)
}
