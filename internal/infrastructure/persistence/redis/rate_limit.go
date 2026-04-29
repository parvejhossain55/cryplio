package redis

import (
	"context"
	"time"
)

// RateLimiter tracks request quotas keyed by subject such as IP or user ID.
type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}

// NoopRateLimiter is a placeholder until Redis-backed rate limiting is implemented.
type NoopRateLimiter struct{}

func NewRateLimiter() *NoopRateLimiter {
	return &NoopRateLimiter{}
}

func (l *NoopRateLimiter) Allow(context.Context, string, int, time.Duration) (bool, error) {
	return true, nil
}
