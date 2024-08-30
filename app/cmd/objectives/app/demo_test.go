package app

import (
	"context"
	"encoding/json"
	"fmt"
	"logbook/models"
	"logbook/models/columns"
	"os"
)

func createRock(ctx context.Context, a *App, uid columns.UserId) (columns.ObjectiveId, error) {
	err := a.RockCreate(ctx, uid)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("RockCreate: %w", err)
	}
	rock, err := a.RockGet(ctx, uid)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("RockGet: %w", err)
	}
	return rock, nil
}

func registerObjectives(ctx context.Context, a *App, uid columns.UserId, parent columns.ObjectiveId, n testfilenode) (columns.ObjectiveId, error) {
	vid, err := a.GetActiveVersion(ctx, parent)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("GetActiveVersion: %w", err)
	}
	registered, err := a.CreateSubtask(ctx, CreateSubtaskParams{
		Creator: uid,
		Parent:  models.Ovid{parent, vid},
		Content: n.Content,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("CreateSubtask: %w", err)
	}
	fmt.Printf("registered %s (%q) on %s %s\n", registered, n.Content, parent, vid)
	for i := 0; i < len(n.Children); i++ {
		_, err := registerObjectives(ctx, a, uid, registered, n.Children[i])
		if err != nil {
			return columns.ZeroObjectId, fmt.Errorf("register(%s/%d): %w", parent, i, err)
		}
	}
	return registered, nil
}

func loadDemo(ctx context.Context, a *App, uid columns.UserId) (columns.ObjectiveId, error) {
	rock, err := createRock(ctx, a, uid)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("createRock: %w", err)
	}
	testfile := []testfilenode{}
	// reading testdata file
	f, err := os.Open("testdata/company.md.json")
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("opening: %w", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&testfile)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("decoding: %w", err)
	}
	for i, n := range testfile {
		_, err := registerObjectives(ctx, a, uid, rock, n)
		if err != nil {
			return columns.ZeroObjectId, fmt.Errorf("register(rock/%d): %w", i, err)
		}
	}
	return rock, nil
}
