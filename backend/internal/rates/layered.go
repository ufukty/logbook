package rates

import (
	"context"
	"fmt"
	"logbook/internal/utils/gsync"
	"logbook/models/columns"

	"golang.org/x/time/rate"
)

var (
	ErrLimitReachedGlobal = fmt.Errorf("binary-wide rate limit has been exceeded")
	ErrLimitReachedUser   = fmt.Errorf("user-specific rate limit has been exceeded")
)

type LimiterParams struct {
	PerSecond float64
	Burst     int
}

type Layered struct {
	binary  *rate.Limiter
	users   gsync.Map[columns.UserId, *rate.Limiter]
	perUser LimiterParams
}

func NewLayered(binaryWide, perUser LimiterParams) *Layered {
	return &Layered{
		binary:  rate.NewLimiter(rate.Limit(binaryWide.PerSecond), binaryWide.Burst),
		perUser: perUser,
	}
}

func (m *Layered) Allow(ctx context.Context, uid columns.UserId) error {
	if !m.binary.Allow() {
		return ErrLimitReachedGlobal
	}

	n := rate.NewLimiter(rate.Limit(m.perUser.PerSecond), m.perUser.Burst)
	user, _ := m.users.LoadOrStore(uid, n)
	if !user.Allow() {
		return ErrLimitReachedUser
	}

	return nil
}

func (m *Layered) Wait(ctx context.Context, uid columns.UserId) error {
	err := m.binary.Wait(ctx)
	if err != nil {
		return fmt.Errorf("checking binary rate limiter: %w", err)
	}

	n := rate.NewLimiter(rate.Limit(m.perUser.PerSecond), m.perUser.Burst)
	user, _ := m.users.LoadOrStore(uid, n)
	err = user.Wait(ctx)
	if err != nil {
		return fmt.Errorf("checking per-user rate limiter: %w", err)
	}

	return nil
}
