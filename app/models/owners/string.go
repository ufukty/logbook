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
	return fmt.Sprintf("%s%s%s:%s (%s)\n",
		strw.Fill("  ", ow.Depth),
		ternary(ow.Folded, "+ ", ""),
		ow.Oid,
		ow.Vid,
		ow.ObjectiveType,
	)
}

}
