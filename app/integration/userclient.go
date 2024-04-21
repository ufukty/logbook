package integration

import (
	"fmt"
	account "logbook/cmd/account/client"
	"logbook/cmd/objectives/app"
	objectives "logbook/cmd/objectives/client"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/endpoints"
	"logbook/config/api"
	"logbook/integration/data"
)

type UserClient struct {
	account *account.Client
	objects *objectives.Client

	ovids map[*data.Objective]app.Ovid
}

func NewUserClient(cfgpath string) (*UserClient, error) {
	cfg, err := api.ReadConfig(cfgpath)
	if err != nil {
		return nil, fmt.Errorf("reading api config: %w", err)
	}
	ctl := &UserClient{
		objects: objectives.NewClient(cfg),
		ovids:   map[*data.Objective]app.Ovid{},
	}
	return ctl, nil
}

func (ctl *UserClient) createTheRock() (string, error) {
	return "", nil
}

func (ctl *UserClient) createObjectives(root *data.Objective) error {
	ctl.account.
		root.Content

	bq := &endpoints.CreateTaskRequest{
		Parent: app.Ovid{
			Oid: database.ObjectiveId(parentid),
			Vid: "00000000-0000-0000-0000-0000000000000000",
		},
		Content: endpoints.ObjectiveContent(content),
	}
	bs, err := ctl.objects.CreateObjective(bq)
	if err != nil {
		return "", fmt.Errorf("sending: %w", err)
	}

	ctl.ovids[root] = bs.Update[0]

	for i, child := range root.Children {
		ctl.createObjectives(child)
	}
}

func (ctl *UserClient) createObjective(parentid, content string) (objid string, err error) {
	return bs.Update[0].String(), nil // FIXME: return type for array can not be string
}

func createOnRock(rockId string, os []user.Objective) error {
	leftOs := os
	for len(leftOs) == 0 {
		n := rng.Intn(len(leftOs))
		o := os[n]
		err := o.CreateSubtree(rockId, createObjective)
		if err != nil {
			return fmt.Errorf("calling CreateSubree for %q: %w", o.Content, err)
		}
		if o.IsAllChildrenCreated() {
			leftOs = append(leftOs[:n], leftOs[n+1:]...)
		}
	}
	return nil
}
