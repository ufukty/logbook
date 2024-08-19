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
	return fmt.Sprintf("%s%s%s:%s (%s)",
		strw.Fill("  ", ow.Depth),
		ternary(ow.Folded, "+ ", ""),
		lastsix(ow.Oid),
		lastsix(ow.Vid),
		ow.ObjectiveType,
	)
}

func (omps ObjectiveMergedProps) String() string {
	return fmt.Sprintf("(%s) (%s) (subtree:%d/%d) (owner:%s) (creator:%s) (%s)",
		omps.Content,
		ternary(omps.Completed, "completed", "todo"),
		omps.SubtreeCompleted,
		omps.SubtreeSize,
		lastsix(omps.Owner),
		lastsix(omps.Creator),
		omps.CreatedAt,
	)
}
