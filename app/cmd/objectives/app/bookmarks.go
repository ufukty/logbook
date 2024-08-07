package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
	"logbook/models/owners"
)

type AddBookmarkParams struct {
	Actor        columns.UserId
	Subject      models.Ovid
	BookmarkName string
}

func (a *App) AddBookmark(ctx context.Context, params AddBookmarkParams) error {
	_, err := a.queries.InsertBookmark(ctx, database.InsertBookmarkParams{
		Uid:    params.Actor,
		Oid:    params.Subject.Oid,
		Vid:    params.Subject.Vid,
		Title:  params.BookmarkName,
		IsRock: false,
	})
	if err != nil {
		return fmt.Errorf("inserting bookmark: %w", err)
	}

	return nil
}

type ListBookmarksParams struct {
	Viewer columns.UserId
}

func (a *App) ListBookmarks(ctx context.Context, params ListBookmarksParams) ([]owners.Bookmark, error) {
	bs, err := a.queries.SelectBookmarks(ctx, params.Viewer)
	if err != nil {
		return nil, fmt.Errorf("selecting bookmarks: %w", err)
	}

	ownbs := []owners.Bookmark{}
	for _, b := range bs {
		ownbs = append(ownbs, owners.Bookmark{
			Bid:       b.Bid,
			Title:     b.Title,
			Oid:       b.Oid,
			Vid:       b.Vid,
			IsRock:    b.IsRock,
			CreatedAt: b.CreatedAt.Time,
		})
	}

	return ownbs, nil
}
