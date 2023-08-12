package contract

import (
	"context"
	"time"
)

type UserRepository interface {
	IncrementUsage(ctx context.Context, userID string, usage int64) (err error)
	DecrementUsage(ctx context.Context, userID string, usage int64) (err error)
	GetUserRateLimit(ctx context.Context, userID string) (limit int, err error)
	ResetUsage(ctx context.Context, end time.Time) (err error)
}
