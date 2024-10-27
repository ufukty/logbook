package app

import (
	"context"
	"fmt"
	"logbook/cmd/groups/database"
	"logbook/models/columns"
	"slices"
)

type CheckMembershipParams struct {
	Uid      columns.UserId
	Gid      columns.GroupId
	Eventual bool
}

func (a *App) getTheGroupsUserIsJoinedInto(ctx context.Context, uid columns.UserId) (map[columns.GroupId]any, error) {
	groups, err := a.oneshot.SelectGroupsByUserTypeMember(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("oneshot.SelectGroupsByUserTypeMember: %w", err)
	}
	acceptance := map[columns.GroupId]any{}
	for _, group := range groups {
		acceptance[group.Gid] = nil
	}
	return acceptance, nil
}

// recurs from the group to its members until
// finds a match between a group's members and,
// the groups the user directly a member of
func (a *App) CheckMembership(ctx context.Context, params CheckMembershipParams) (bool, error) {
	joined, err := a.getTheGroupsUserIsJoinedInto(ctx, params.Uid)
	if err != nil {
		return false, fmt.Errorf("a.getTheGroupsUserIsJoinedInto: %w", err)
	}

	// direct
	if _, ok := joined[params.Gid]; ok {
		return true, nil
	}

	if !params.Eventual {
		return false, nil
	}

	// indirect
	queue := []columns.GroupId{params.Gid}
	checked := map[columns.GroupId]any{}
	for _, g := range queue {
		checked[g] = nil

		members, err := a.oneshot.SelectUserTypeGroupMembers(ctx, g)
		if err != nil {
			return false, fmt.Errorf("oneshot.SelectUserTypeGroupMembers: %w", err)
		}

		if slices.IndexFunc(members, func(member database.GroupMemberUser) bool {
			return member.Member == params.Uid
		}) != -1 {
			return true, nil
		}

		subgroups, err := a.oneshot.SelectGroupTypeGroupMembers(ctx, g)
		if err != nil {
			return false, fmt.Errorf("oneshot.SelectGroupTypeGroupMembers: %w", err)
		}

		for _, subgroup := range subgroups {
			if _, ok := joined[subgroup.Gid]; ok {
				return true, nil
			}

			_, ok := checked[subgroup.Gid]
			if !ok {
				queue = append(queue, subgroup.Gid)
			}
		}
	}

	return false, nil
}
