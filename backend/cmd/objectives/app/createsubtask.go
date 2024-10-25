package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/app/l2"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
)

const LimitSubitems = 20

var ErrTooManySubitems = fmt.Errorf("limit for subitems per objective has been reached on the parent")

type CreateSubtaskParams struct {
	Creator columns.UserId
	Parent  models.Ovid
	Content string
}

// TODO: check privileges on parent
// DONE: create operations
// DONE: props
// DONE: transaction-commit-rollback
// DONE: bubblink
// DONE: mark active version for promoted ascendants
// DONE: enforce 20-subtask limit on the parent
// TODO: apply auto-merge on non-conflixcting concurrent updates
// TODO: apply detaching from active version on conflicting concurrent updates
// TODO: invalidate the unfolded subtree size cache for each ascendant and each viewer & trigger recalculation
func (a *App) CreateSubtask(ctx context.Context, params CreateSubtaskParams) (models.Ovid, error) {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	activepath, err := a.l2.ListActivePathToRock(ctx, q, params.Parent)
	if err == l2.ErrLeftBehind {
		return models.ZeroOvid, l2.ErrLeftBehind
	} else if err != nil {
		return models.ZeroOvid, fmt.Errorf("listActivePathToRock: %w", err)
	}

	pObj, err := q.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: params.Parent.Oid,
		Vid: params.Parent.Vid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("SelectObjective: %w", err)
	}

	pBups, err := q.SelectBottomUpProps(ctx, pObj.Bupid)
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("SelectBottomUpProps/parent: %w", err)
	}

	if pBups.Children >= LimitSubitems {
		return models.ZeroOvid, ErrTooManySubitems
	}

	op, err := q.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.Parent.Oid,
		Subjectvid: params.Parent.Vid,
		Actor:      params.Creator,
		OpType:     database.OpTypeObjInit,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertOperation/child: %w", err)
	}

	_, err = q.InsertOpObjInit(ctx, database.InsertOpObjInitParams{
		Opid:    op.Opid,
		Content: params.Content,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertOpObjInit: %w", err)
	}

	props, err := q.InsertProperties(ctx, database.InsertPropertiesParams{
		Content:   params.Content,
		Completed: false,
		Creator:   params.Creator,
		Owner:     params.Creator,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertProperties/child: %w", err)
	}

	bup, err := q.InsertBottomUpProps(ctx, database.InsertBottomUpPropsParams{
		Children:         0,
		SubtreeSize:      0,
		SubtreeCompleted: 0,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertBottomUpProps/child: %w", err)
	}

	obj, err := q.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: op.Opid,
		Pid:       props.Pid,
		Bupid:     bup.Bupid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertNewObjective: %w", err)
	}

	_, err = q.InsertActiveVidForObjective(ctx, database.InsertActiveVidForObjectiveParams{
		Oid: obj.Oid,
		Vid: obj.Vid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertActiveVidForObjective: %w", err)
	}

	pOp, err := q.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: pObj.Oid,
		Subjectvid: pObj.Vid,
		Actor:      params.Creator,
		OpType:     database.OpTypeObjCreateSubtask,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertOperation: %w", err)
	}

	_, err = q.InsertOpObjCreateSubtask(ctx, database.InsertOpObjCreateSubtaskParams{
		Opid: pOp.Opid,
		Soid: obj.Oid,
		Svid: obj.Vid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertOpObjCreateSubtask: %w", err)
	}

	pBups.Children += 1
	pBups.SubtreeSize += 1
	pBupsUpd, err := q.InsertBottomUpProps(ctx, database.InsertBottomUpPropsParams{
		Children:         pBups.Children,
		SubtreeSize:      pBups.SubtreeSize,
		SubtreeCompleted: pBups.SubtreeCompleted,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertBottomUpProps/parent: %w", err)
	}

	pObjUpd, err := q.InsertUpdatedObjective(ctx, database.InsertUpdatedObjectiveParams{
		Oid:       pObj.Oid,
		Based:     pObj.Vid,
		CreatedBy: pOp.Opid,
		Pid:       pObj.Pid,
		Bupid:     pBupsUpd.Bupid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertUpdatedObjective: %w", err)
	}

	_, err = q.InsertNewLink(ctx, database.InsertNewLinkParams{
		SupOid: pObjUpd.Oid,
		SupVid: pObjUpd.Vid,
		SubOid: obj.Oid,
		SubVid: obj.Vid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("InsertNewLink: %w", err)
	}

	sublinks, err := q.SelectSubLinks(ctx, database.SelectSubLinksParams{
		SupOid: pObj.Oid,
		SupVid: pObj.Vid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("SelectSubLinks: %w", err)
	}
	for _, sublink := range sublinks {
		_, err := q.InsertUpdatedLink(ctx, database.InsertUpdatedLinkParams{
			SupOid:            pObjUpd.Oid,
			SupVid:            pObjUpd.Vid,
			SubOid:            sublink.SubOid,
			SubVid:            sublink.SubVid,
			CreatedAtOriginal: sublink.CreatedAtOriginal,
		})
		if err != nil {
			return models.ZeroOvid, fmt.Errorf("InsertUpdatedLink: %w", err)
		}
	}

	_, err = q.UpdateActiveVidForObjective(ctx, database.UpdateActiveVidForObjectiveParams{
		Oid: pObj.Oid,
		Vid: pObjUpd.Vid,
	})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("UpdateActiveVidForObjective: %w", err)
	}

	activepath[0].Vid = pObjUpd.Vid
	_, err = a.bubblink(ctx, q, activepath, pOp, bubblinkDeltaValues{SubtreeSize: 1})
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("bubblink: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.ZeroOvid, fmt.Errorf("commit: %w", err)
	}

	return models.Ovid{Oid: obj.Oid, Vid: obj.Vid}, nil
}
