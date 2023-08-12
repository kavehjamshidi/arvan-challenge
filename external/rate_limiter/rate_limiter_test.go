package rate_limiter

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type rateLimiterSuite struct {
	suite.Suite
	redisClient *redis.Client
	mock        redismock.ClientMock
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, &rateLimiterSuite{})
}

func (r *rateLimiterSuite) SetupSuite() {
	db, mock := redismock.NewClientMock()
	r.redisClient = db
	r.mock = mock
}

func (r *rateLimiterSuite) TestGetUserLimit() {
	r.Run("success", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectGet("quota:123456").SetVal("200")

		quota, err := limiter.GetUserLimit(context.TODO(), "123456")

		r.NoError(err)
		r.Equal(200, quota)
	})

	r.Run("not found", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectGet("quota:123456").RedisNil()

		quota, err := limiter.GetUserLimit(context.TODO(), "123456")

		r.NoError(err)
		r.Equal(0, quota)
	})

	r.Run("fail - redis error", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectGet("quota:123456").SetErr(errors.New("unknown error"))

		quota, err := limiter.GetUserLimit(context.TODO(), "123456")

		r.ErrorContains(err, "unknown error")
		r.Equal(0, quota)
	})

	r.Run("fail - quota not integer", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectGet("quota:123456").SetVal("abc")

		quota, err := limiter.GetUserLimit(context.TODO(), "123456")

		r.ErrorContains(err, "strconv.Atoi")
		r.Equal(0, quota)
	})
}

func (r *rateLimiterSuite) TestSetUserLimit() {
	r.Run("success", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectSet("quota:123456", 1000, 0).SetVal("1")

		err := limiter.SetUserLimit(context.TODO(), "123456", 1000)

		r.NoError(err)
	})

	r.Run("failed", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectSet("quota:123456", 1000, 0).SetErr(errors.New("unknown error"))

		err := limiter.SetUserLimit(context.TODO(), "123456", 1000)

		r.ErrorContains(err, "unknown error")
	})
}

func (r *rateLimiterSuite) TestIsAllowed() {
	now := time.Now().UnixNano()
	windowSize := windowsSizeSeconds
	luaScript := `
  local key = KEYS[1]
  local now = tonumber(ARGV[1])
  local window = tonumber(ARGV[2])
  local windowsNano = window * 1000000000
  local limit = tonumber(ARGV[3])
  local clearBefore = now - windowsNano
  redis.call('ZREMRANGEBYSCORE', key, 0, clearBefore)
  local amount = redis.call('ZCARD', key)
  if amount <= limit then
  redis.call('ZADD', key, now, now)
  end
  redis.call('EXPIRE', key, window)
  return limit - amount
`

	r.Run("success - allowed", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectEval(luaScript, []string{"rate:123456"}, now, windowSize, 100).SetVal(int64(10))

		isAllowed, err := limiter.IsAllowed(context.TODO(), "123456", 100)

		r.NoError(err)
		r.True(isAllowed)
	})

	r.Run("success - not allowed", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectEval(luaScript, []string{"rate:123456"}, now, windowSize, 100).SetVal(int64(-1))

		isAllowed, err := limiter.IsAllowed(context.TODO(), "123456", 100)

		r.NoError(err)
		r.False(isAllowed)
	})

	r.Run("failed - error", func() {
		limiter := NewRedisRateLimiter(r.redisClient)

		r.mock.ExpectEval(luaScript, []string{"rate:123456"}, now, windowSize, 100).SetErr(errors.New("unknown error"))

		isAllowed, err := limiter.IsAllowed(context.TODO(), "123456", 100)

		r.ErrorContains(err, "unknown error")
		r.False(isAllowed)
	})
}
