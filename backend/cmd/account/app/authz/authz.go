package authz

import (
	"fmt"
	"logbook/cmd/account/database"
)

var (
	NoAuthorization = fmt.Errorf("no authorization")
	UnderAuthorized = fmt.Errorf("under authorized")
)

type Authorization struct {
	queries *database.Queries
}

func New(q *database.Queries) *Authorization {
	return &Authorization{
		queries: q,
	}
}
