package integration

import (
	"fmt"
)

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
