package transports

import (
	"logbook/models/columns"
)

type (
	GroupIds []columns.GroupId
	UserIds  []columns.UserId
)

func (gs GroupIds) Validate() any {
	issues := []any{}
	for _, g := range gs {
		if issue := g.Validate(); issue != nil {
			issues = append(issues, issue)
		}
	}
	if len(issues) > 0 {
		return issues
	}
	return nil
}
