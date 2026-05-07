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
	GetByUsernameWithStats(ctx context.Context, username string) (*User, *UserStats, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, limit, offset int) ([]User, error)
	CountUsers(ctx context.Context) (int, error)
	IncrementLogin(ctx context.Context, userID uuid.UUID) error
	UpdateLastSeen(ctx context.Context, userID uuid.UUID) error
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error)
	GetStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
	// OAuth
	GetOAuthByProviderID(ctx context.Context, provider, providerUserID string) (*UserOAuth, error)
	CreateOAuth(ctx context.Context, oauth *UserOAuth) error
	UpdateOAuth(ctx context.Context, oauth *UserOAuth) error
	GetOAuthByUserID(ctx context.Context, userID uuid.UUID) ([]UserOAuth, error)
	DeleteOAuth(ctx context.Context, id uuid.UUID) error
	// Email Verification
	CreateEmailVerificationToken(ctx context.Context, token *EmailVerificationToken) error
	GetEmailVerificationTokenByHash(ctx context.Context, tokenHash string) (*EmailVerificationToken, error)
	GetEmailVerificationTokenByUserID(ctx context.Context, userID uuid.UUID) (*EmailVerificationToken, error)
	MarkEmailVerificationTokenVerified(ctx context.Context, id int) error
	// Password Reset
	CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	GetPasswordResetTokenByUserID(ctx context.Context, userID uuid.UUID) (*PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, id int) error
	// Two-Factor Pending
	CreateTwoFactorPending(ctx context.Context, pending *TwoFactorPending) error
	GetTwoFactorPendingByUserID(ctx context.Context, userID uuid.UUID) (*TwoFactorPending, error)
	DeleteTwoFactorPending(ctx context.Context, userID uuid.UUID) error
	// Sessions
	CreateSession(ctx context.Context, session *UserSession) error
	GetSession(ctx context.Context, tokenID string) (*UserSession, error)
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error)
	DeleteSession(ctx context.Context, tokenID string) error
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
	UpdateSessionLastUsed(ctx context.Context, tokenID string) error

	// User Payment Methods
	CreateUserPaymentMethod(ctx context.Context, pm *UserPaymentMethod) error
	GetUserPaymentMethod(ctx context.Context, pmID uuid.UUID) (*UserPaymentMethod, error)
	GetUserPaymentMethods(ctx context.Context, userID uuid.UUID) ([]UserPaymentMethod, error)
	UpdateUserPaymentMethod(ctx context.Context, pm *UserPaymentMethod) error
	DeleteUserPaymentMethod(ctx context.Context, pmID uuid.UUID) error
	SetDefaultUserPaymentMethod(ctx context.Context, userID, pmID uuid.UUID) error
}
