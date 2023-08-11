package contract

import "context"

type RateLimitService interface {
	CheckRateLimit(ctx context.Context, userID string) (err error)
}
