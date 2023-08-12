package rate_limit

import (
	"context"
	contract2 "github.com/kavehjamshidi/arvan-challenge/external/rate_limiter/contract"
	"github.com/kavehjamshidi/arvan-challenge/repository/user/contract"
	contract3 "github.com/kavehjamshidi/arvan-challenge/service/rate_limit/contract"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
)

type rateLimitService struct {
	userRepo    contract.UserRepository
	rateLimiter contract2.RateLimiter
}

func NewRateLimitService(
	userRepo contract.UserRepository,
	rateLimiter contract2.RateLimiter,
) contract3.RateLimitService {
	return rateLimitService{
		userRepo:    userRepo,
		rateLimiter: rateLimiter,
	}
}

func (rs rateLimitService) CheckRateLimit(ctx context.Context, userID string) error {
	limit, err := rs.rateLimiter.GetUserLimit(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "CheckRateLimit")
	}

	if limit == 0 {
		limit, err = rs.userRepo.GetUserRateLimit(ctx, userID)
		if err != nil {
			return errors.Wrap(err, "CheckRateLimit")
		}

		if limit == 0 {
			return errors.Wrap(utils.ErrUserRateLimitNotSet, "CheckRateLimit")
		}

		rs.rateLimiter.SetUserLimit(ctx, userID, limit)
	}

	allowed, err := rs.rateLimiter.IsAllowed(ctx, userID, limit)
	if err != nil {
		return errors.Wrap(err, "CheckRateLimit")
	}

	if !allowed {
		return errors.Wrap(utils.ErrTooManyRequests, "CheckRateLimit")
	}

	return nil
}
