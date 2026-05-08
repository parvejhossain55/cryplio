package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// NewClient creates and pings a go-redis client.
// addr is in "host:port" format (e.g. "localhost:6379").
func NewClient(addr, password string, db int) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect to redis at %s: %w", addr, err)
	}
	return client, nil
}
