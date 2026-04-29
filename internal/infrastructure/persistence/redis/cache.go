package redis

import (
	"context"
	"time"
)

// Cache provides a narrow cache abstraction for infrastructure adapters.
type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

// NoopCache is a placeholder until Redis-backed caching is implemented.
type NoopCache struct{}

func NewCache() *NoopCache {
	return &NoopCache{}
}

func (c *NoopCache) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}

func (c *NoopCache) Get(context.Context, string) ([]byte, error) {
	return nil, nil
}

func (c *NoopCache) Delete(context.Context, string) error {
	return nil
}
