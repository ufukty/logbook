package app

import (
	"context"
	"fmt"
	"testing"
)

func TestListChildrenUpToDateResults(t *testing.T) {
	a, uid, err := testdeps()
	if err != nil {
		t.Fatal(fmt.Errorf("prep, testdeps: %w", err))
	}
	defer a.pool.Close()

	rock, err := loadDemo(context.Background(), a, uid, false)
	if err != nil {
		t.Fatal(fmt.Errorf("prep, CreateDemoFileInDfsOrder: %w", err))
	}
	fmt.Println("rock:", rock)
	children, err := a.ListChildren(context.Background(), rock)
	if err != nil {
		t.Fatal(fmt.Errorf("ListChildren: %w", err))
	}

	for _, child := range children {
		got, err := a.GetActiveVersion(context.Background(), child.Oid)
		if err != nil {
			t.Fatal(fmt.Errorf("GetActiveVersion(%s): %w", child.Oid, err))
		}
		if got != child.Vid {
			t.Errorf("assert, outdated version returned by GetActiveVersion(%s): want %s, got %s", child.Oid, child.Vid, got)
		}
	}
}
