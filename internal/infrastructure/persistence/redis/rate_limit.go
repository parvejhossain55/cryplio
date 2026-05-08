package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// RateLimiter tracks request quotas keyed by subject such as IP or user ID.
type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}

// RedisRateLimiter implements a sliding-window rate limiter backed by a Redis
// sorted set. Each request is recorded as a member with score = Unix nanosecond
// timestamp. Old entries outside the window are pruned on every check.
//
// This implementation is safe for horizontally-scaled deployments.
type RedisRateLimiter struct {
	client *goredis.Client
}

// NewRedisRateLimiter creates a rate limiter backed by the given Redis client.
func NewRedisRateLimiter(client *goredis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{client: client}
}

// Allow returns true if the key has not exceeded limit requests within window.
// It uses a Lua script to atomically prune old entries and check + increment
// the counter in a single round-trip.
var slidingWindowScript = goredis.NewScript(`
local key    = KEYS[1]
local now    = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit  = tonumber(ARGV[3])
local cutoff = now - window

redis.call("ZREMRANGEBYSCORE", key, "-inf", cutoff)
local count = redis.call("ZCARD", key)
if count >= limit then
    return 0
end
redis.call("ZADD", key, now, now)
redis.call("PEXPIRE", key, window)
return 1
`)

func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now().UnixMilli()
	windowMs := window.Milliseconds()

	result, err := slidingWindowScript.Run(ctx, r.client,
		[]string{fmt.Sprintf("rl:%s", key)},
		now, windowMs, limit,
	).Int()
	if err != nil {
		// On Redis error, fail open (allow the request) to avoid blocking all traffic.
		return true, fmt.Errorf("rate limiter redis error: %w", err)
	}
	return result == 1, nil
}

// NoopRateLimiter always allows every request. Used as a placeholder.
type NoopRateLimiter struct{}

func NewNoopRateLimiter() *NoopRateLimiter { return &NoopRateLimiter{} }

func (l *NoopRateLimiter) Allow(context.Context, string, int, time.Duration) (bool, error) {
	return true, nil
}
