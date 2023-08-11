package contract

import "context"

type RateLimiter interface {
	GetUserLimit(ctx context.Context, userID string) (limit int, err error)
	SetUserLimit(ctx context.Context, userID string, limit int) (err error)
	IsAllowed(ctx context.Context, userID string, limit int) (isAllowed bool, err error)
}
