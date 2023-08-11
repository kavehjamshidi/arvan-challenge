package rate_limiter

import (
	"context"
	"fmt"
	"github.com/kavehjamshidi/arvan-challenge/external/rate_limiter/contract"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

const (
	ratePrefix  = "rate"
	quotaPrefix = "quota"

	windowsSizeSeconds = 60 // One minute
)

type rateLimiter struct {
	redisClient *redis.Client
}

func NewRedisRateLimiter(redisClient *redis.Client) contract.RateLimiter {
	return rateLimiter{redisClient: redisClient}
}

func (r rateLimiter) SetUserLimit(ctx context.Context, userID string, limit int) error {
	key := r.generateKey(userID, quotaPrefix)

	err := r.redisClient.Set(ctx, key, limit, 0).Err()
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "SetUserLimit")
	}

	return nil
}

func (r rateLimiter) GetUserLimit(ctx context.Context, userID string) (int, error) {
	key := r.generateKey(userID, quotaPrefix)

	res, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Println(err)
		return 0, errors.Wrap(err, "GetUserQuota")
	}

	quota, err := strconv.Atoi(res)
	if err != nil {
		log.Println(err)
		return 0, errors.Wrap(err, "GetUserQuota")
	}

	return quota, nil
}

func (r rateLimiter) IsAllowed(ctx context.Context, userID string, limit int) (bool, error) {
	now := time.Now().UnixNano()
	windowSize := windowsSizeSeconds
	key := r.generateKey(userID, ratePrefix)

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
	res, err := r.redisClient.Eval(
		ctx,
		luaScript,     // script
		[]string{key}, // KEYS
		now,           // ARGV[1]
		windowSize,    // ARGV[2]
		limit,         // ARGV[3]
	).Result()
	if err != nil {
		log.Println(err)
		return false, errors.Wrap(err, "Allowed")
	}

	remaining, _ := res.(int64)
	if remaining < 0 {
		return false, nil
	}

	return true, nil
}

func (r rateLimiter) generateKey(userID string, prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, userID)
}
