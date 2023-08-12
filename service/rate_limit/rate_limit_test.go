package rate_limit

import (
	"context"
	"github.com/kavehjamshidi/arvan-challenge/mocks"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type rateLimitSuite struct {
	suite.Suite
	userRepo    *mocks.UserRepository
	rateLimiter *mocks.RateLimiter
}

func TestRateLimitSuite(t *testing.T) {
	suite.Run(t, &rateLimitSuite{})
}

func (r *rateLimitSuite) SetupSubTest() {
	mockUserRepo := mocks.NewUserRepository(r.T())
	mockRateLimiter := mocks.NewRateLimiter(r.T())

	r.userRepo = mockUserRepo
	r.rateLimiter = mockRateLimiter
}

func (r *rateLimitSuite) TestCheckRateLimit() {
	userID := "123456"

	r.Run("error - failed in getting limit from redis", func() {
		rateLimitSVC := NewRateLimitService(r.userRepo, r.rateLimiter)

		r.rateLimiter.On("GetUserLimit", context.TODO(), userID).
			Return(0, errors.New("unknown error"))

		err := rateLimitSVC.CheckRateLimit(context.TODO(), userID)

		r.ErrorContains(err, "unknown error")
	})

	r.Run("error - failed in getting limit from db", func() {
		rateLimitSVC := NewRateLimitService(r.userRepo, r.rateLimiter)

		r.rateLimiter.On("GetUserLimit", context.TODO(), userID).
			Return(0, nil)

		r.userRepo.On("GetUserRateLimit", context.TODO(), userID).
			Return(0, errors.New("unknown error 2"))

		err := rateLimitSVC.CheckRateLimit(context.TODO(), userID)

		r.ErrorContains(err, "unknown error 2")
	})

	r.Run("error - no limit set", func() {
		rateLimitSVC := NewRateLimitService(r.userRepo, r.rateLimiter)

		r.rateLimiter.On("GetUserLimit", context.TODO(), userID).
			Return(0, nil)

		r.userRepo.On("GetUserRateLimit", context.TODO(), userID).
			Return(0, nil)

		err := rateLimitSVC.CheckRateLimit(context.TODO(), userID)

		r.ErrorIs(err, utils.ErrUserRateLimitNotSet)
	})

	r.Run("error - failed to get allowed", func() {
		rateLimitSVC := NewRateLimitService(r.userRepo, r.rateLimiter)

		r.rateLimiter.On("GetUserLimit", context.TODO(), userID).
			Return(0, nil)

		r.userRepo.On("GetUserRateLimit", context.TODO(), userID).
			Return(100, nil)

		r.rateLimiter.On("SetUserLimit", context.TODO(), userID, 100).
			Return(nil)

		r.rateLimiter.On("IsAllowed", context.TODO(), userID, 100).
			Return(false, errors.New("unknown error"))

		err := rateLimitSVC.CheckRateLimit(context.TODO(), userID)

		r.ErrorContains(err, "unknown error")
	})

	r.Run("error - not allowed", func() {
		rateLimitSVC := NewRateLimitService(r.userRepo, r.rateLimiter)

		r.rateLimiter.On("GetUserLimit", context.TODO(), userID).
			Return(0, nil)

		r.userRepo.On("GetUserRateLimit", context.TODO(), userID).
			Return(100, nil)

		r.rateLimiter.On("SetUserLimit", context.TODO(), userID, 100).
			Return(nil)

		r.rateLimiter.On("IsAllowed", context.TODO(), userID, 100).
			Return(false, nil)

		err := rateLimitSVC.CheckRateLimit(context.TODO(), userID)

		r.ErrorIs(err, utils.ErrTooManyRequests)
	})

	r.Run("success", func() {
		rateLimitSVC := NewRateLimitService(r.userRepo, r.rateLimiter)

		r.rateLimiter.On("GetUserLimit", context.TODO(), userID).
			Return(0, nil)

		r.userRepo.On("GetUserRateLimit", context.TODO(), userID).
			Return(100, nil)

		r.rateLimiter.On("SetUserLimit", context.TODO(), userID, 100).
			Return(nil)

		r.rateLimiter.On("IsAllowed", context.TODO(), userID, 100).
			Return(true, nil)

		err := rateLimitSVC.CheckRateLimit(context.TODO(), userID)

		r.NoError(err)
	})
}
