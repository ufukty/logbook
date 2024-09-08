package app

import (
	"context"
	"encoding/json"
	"fmt"
	"logbook/models"
	"logbook/models/columns"
	"os"
)

func createRock(ctx context.Context, a *App, uid columns.UserId, debug bool) (columns.ObjectiveId, error) {
	err := a.RockCreate(ctx, uid)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("RockCreate: %w", err)
	}
	rock, err := a.RockGet(ctx, uid)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("RockGet: %w", err)
	}
	if debug {
		fmt.Println("rock:", rock)
	}
	return rock, nil
}

func registerObjectives(ctx context.Context, a *App, uid columns.UserId, parent columns.ObjectiveId, n testfilenode, debug bool) (columns.ObjectiveId, error) {
	vid, err := a.GetActiveVersion(ctx, parent)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("GetActiveVersion: %w", err)
	}
	if debug {
		fmt.Printf("registering (%s %s) <- (%s)\n", parent, vid, n.Content)
	}
	registered, err := a.CreateSubtask(ctx, CreateSubtaskParams{
		Creator: uid,
		Parent:  models.Ovid{parent, vid},
		Content: n.Content,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("CreateSubtask: %w", err)
	}
	if debug {
		fmt.Printf("registered: %s\n", registered)
	}
	for i := 0; i < len(n.Children); i++ {
		_, err := registerObjectives(ctx, a, uid, registered.Oid, n.Children[i], debug)
		if err != nil {
			return columns.ZeroObjectId, fmt.Errorf("register(%s/%d): %w", parent, i, err)
		}
	}
	return registered.Oid, nil
}

func loadDemo(ctx context.Context, a *App, uid columns.UserId, debug bool) (models.Ovid, error) {
	var rock models.Ovid
	var err error
	rock.Oid, err = createRock(ctx, a, uid, debug)
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("createRock: %w", err)
	}
	testfile := []testfilenode{}
	// reading testdata file
	f, err := os.Open("testdata/company.md.json")
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("opening: %w", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&testfile)
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("decoding: %w", err)
	}
	for i, n := range testfile {
		_, err := registerObjectives(ctx, a, uid, rock.Oid, n, debug)
		if err != nil {
			return models.ZeroOvid, fmt.Errorf("register(rock/%d): %w", i, err)
		}
	}
	rock.Vid, err = a.GetActiveVersion(context.Background(), rock.Oid)
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("GetActiveVersion: %w", err)
	}
	return rock, nil
}
