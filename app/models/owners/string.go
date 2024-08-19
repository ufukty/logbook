package owners

import (
	"fmt"
	"logbook/internal/utilities/strw"
)

func (ow DocumentItem) String() string {
	fold := ""
	if ow.Folded {
		fold = "+ "
	}
	return fmt.Sprintf("%s%s%s:%s (%s)\n", strw.Fill("  ", ow.Depth), fold, ow.Oid, ow.Vid, ow.ObjectiveType)
}
