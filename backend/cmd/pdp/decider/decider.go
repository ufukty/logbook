package decider

import (
	"fmt"
	groups "logbook/cmd/groups/client"
	objectives "logbook/cmd/objectives/client"
)

var ErrUnderAuthorized = fmt.Errorf("under authorized")

type Decider struct {
	groups     *groups.Client
	objectives *objectives.Client
}

func New(g *groups.Client, o *objectives.Client) *Decider {
	return &Decider{
		groups:     g,
		objectives: o,
	}
}
