package redis

import (
	"context"
	"time"

	identity "cryplio/internal/domain/identity"
)

// SessionStore defines session persistence backed by Redis or a compatible cache.
type SessionStore interface {
	Save(ctx context.Context, session *identity.UserSession, ttl time.Duration) error
	Get(ctx context.Context, tokenID string) (*identity.UserSession, error)
	Delete(ctx context.Context, tokenID string) error
}

// InMemorySessionStore is a compile-safe placeholder until Redis is wired in.
type InMemorySessionStore struct{}

func NewSessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{}
}

func (s *InMemorySessionStore) Save(context.Context, *identity.UserSession, time.Duration) error {
	return nil
}

func (s *InMemorySessionStore) Get(context.Context, string) (*identity.UserSession, error) {
	return nil, nil
}

func (s *InMemorySessionStore) Delete(context.Context, string) error {
	return nil
}
