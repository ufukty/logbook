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

func (ow DocumentItem) String() string {
	return fmt.Sprintf("%s%s%s:%s (%s)",
		strw.Fill("  ", ow.Depth),
		ternary(ow.Folded, "+ ", ""),
		ow.Oid,
		ow.Vid,
		ow.ObjectiveType,
	)
}

func (omps ObjectiveMergedProps) String() string {
	return fmt.Sprintf("(%s) (%s) (subtree:%d/%d) (owner:%s) (creator:%s) (%s)",
		omps.Content,
		ternary(omps.Completed, "completed", "todo"),
		omps.SubtreeCompleted,
		omps.SubtreeSize,
		omps.Owner,
		omps.Creator,
		omps.CreatedAt,
	)
}
