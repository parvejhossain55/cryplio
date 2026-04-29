package identity

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines data access for users
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, limit, offset int) ([]User, error)
	IncrementLogin(ctx context.Context, userID uuid.UUID) error
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error)
	GetStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
	// OAuth
	GetOAuthByProviderID(ctx context.Context, provider, providerUserID string) (*UserOAuth, error)
	CreateOAuth(ctx context.Context, oauth *UserOAuth) error
	UpdateOAuth(ctx context.Context, oauth *UserOAuth) error
	GetOAuthByUserID(ctx context.Context, userID uuid.UUID) ([]UserOAuth, error)
	DeleteOAuth(ctx context.Context, id uuid.UUID) error
}
