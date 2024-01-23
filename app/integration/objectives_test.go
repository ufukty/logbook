package database

import (
	"encoding/json"
	"fmt"
	"logbook/cmd/tasks/database"
	"logbook/cmd/tasks/endpoints"
	"math/rand"
	"os"
)

var rng = rand.New(rand.NewSource(1))

type Objective struct {
	Content  string      `json:"content"`
	Children []Objective `json:"children"`

	IsCreated               bool   `json:"-"`
	NumberOfChildrenCreated int    `json:"-"`
	Id                      string `json:"-"` // the Id returned by App when it is created
}

func (a *Objective) IsAllChildrenCreated() bool {
	return len(a.Children) == a.NumberOfChildrenCreated
}

func (a *Objective) ChildrenThatHasUncreatedChildren() []Objective {
	children := []Objective{}
	for _, child := range a.Children {
		if !child.IsAllChildrenCreated() {
			children = append(children, child)
		}
	}
	return children
}

func (a *Objective) CreateSubtree(parentId string, creator func(parentId, content string) (string, error)) error {
	if !a.IsCreated {
		id, err := creator(parentId, a.Content)
		if err != nil {
			return fmt.Errorf("calling creator from %q: %w", a.Content, err)
		}
		a.Id = id
		a.IsCreated = true
		return nil
	} else {
		if a.IsAllChildrenCreated() {
			return nil
		}
		children := a.ChildrenThatHasUncreatedChildren()
		child := children[rng.Intn(len(children))]
		err := child.CreateSubtree(a.Id, creator)
		if err != nil {
			return fmt.Errorf("calling CreateSubtree on child %q: %w", a.Content, err)
		}
		a.NumberOfChildrenCreated++
		return nil
	}
}

func (o Objective) CountSubtree() int {
	pop := 1
	for _, c := range o.Children {
		pop += c.CountSubtree()
	}
	return pop
}

func load() ([]Objective, error) {
	f, err := os.Open("testdata/company.json")
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}
	defer f.Close()

	os := &[]Objective{}
	err = json.NewDecoder(f).Decode(os)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return *os, nil
}

func createTheRock() (string, error) {
	return "", nil
}

func createObjective(parentid, content string) (objid string, err error) {
	bq := endpoints.CreateTaskRequest{ParentOid: database.ObjectiveId(parentid), Text: content}
	bs, err := bq.Send()
	if err != nil {
		return "", fmt.Errorf("sending: %w", err)
	}
	return string(bs.Created.Oid), nil
}

func createOnRock(rockId string, os []Objective) error {
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
