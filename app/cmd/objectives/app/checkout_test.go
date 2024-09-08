package app

import (
	"context"
	"fmt"
	"testing"
)

func TestHistoryAndCheckout(t *testing.T) {
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
	first := children[0]
	fmt.Println("child:", first)

	history, err := a.GetObjectiveHistory(context.Background(), GetObjectiveHistoryParams{
		Subject:               first,
		IncludeAdministrative: false,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act, GetObjectiveHistory: %w", err))
	}

	targetVid := history[len(history)-1].Version
	fmt.Println("target vid:", targetVid)
	err = a.Checkout(context.Background(), CheckoutParams{
		User:    uid,
		Subject: first,
		To:      targetVid,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act, Checkout: %w", err))
	}

	active, err := a.GetActiveVersion(context.Background(), first.Oid)
	if err != nil {
		t.Fatal(fmt.Errorf("assert, GetActiveVersion: %w", err))
	}

	if active != targetVid {
		t.Errorf("assert, active version is not applied. want %s, got %s", targetVid, active)
	}
}
