package transports

import (
	"logbook/models/columns"
)

type (
	GroupIds []columns.GroupId
	UserIds  []columns.UserId
)

func collect[E interface{ Validate() any }](s []E) any {
	issues := []any{}
	for _, e := range s {
		if issue := e.Validate(); issue != nil {
			issues = append(issues, issue)
		}
	}
	if len(issues) > 0 {
		return issues
	}
	return nil
}

func (gs GroupIds) Validate() any { return collect(gs) }
func (us UserIds) Validate() any  { return collect(us) }
