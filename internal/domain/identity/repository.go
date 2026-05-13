package identity

import (
	"context"

	"github.com/google/uuid"
)

// ─── Segregated repository interfaces ────────────────────────────────────────
// Each interface covers a single concern. Consumers (services, mocks, tests)
// should depend only on the interface they actually need.

// UserCoreRepository handles basic user CRUD and login-tracking queries.
type UserCoreRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByUsernameWithStats(ctx context.Context, username string) (*User, *UserStats, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, limit, offset int, searchQuery string, status string) ([]User, int, error)
	CountUsers(ctx context.Context) (int, error)
	IncrementLogin(ctx context.Context, userID uuid.UUID) error
	UpdateLastSeen(ctx context.Context, userID uuid.UUID) error
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error)
	GetStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
}

// OAuthRepository manages third-party OAuth provider links.
type OAuthRepository interface {
	GetOAuthByProviderID(ctx context.Context, provider, providerUserID string) (*UserOAuth, error)
	CreateOAuth(ctx context.Context, oauth *UserOAuth) error
	UpdateOAuth(ctx context.Context, oauth *UserOAuth) error
	GetOAuthByUserID(ctx context.Context, userID uuid.UUID) ([]UserOAuth, error)
	DeleteOAuth(ctx context.Context, id uuid.UUID) error
}

// EmailVerificationRepository manages email-verification tokens.
type EmailVerificationRepository interface {
	CreateEmailVerificationToken(ctx context.Context, token *EmailVerificationToken) error
	GetEmailVerificationTokenByHash(ctx context.Context, tokenHash string) (*EmailVerificationToken, error)
	GetEmailVerificationTokenByUserID(ctx context.Context, userID uuid.UUID) (*EmailVerificationToken, error)
	MarkEmailVerificationTokenVerified(ctx context.Context, id int) error
}

// PasswordResetRepository manages password-reset tokens.
type PasswordResetRepository interface {
	CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	GetPasswordResetTokenByUserID(ctx context.Context, userID uuid.UUID) (*PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, id int) error
}

// TwoFactorRepository manages in-progress 2FA setup state.
type TwoFactorRepository interface {
	CreateTwoFactorPending(ctx context.Context, pending *TwoFactorPending) error
	GetTwoFactorPendingByUserID(ctx context.Context, userID uuid.UUID) (*TwoFactorPending, error)
	DeleteTwoFactorPending(ctx context.Context, userID uuid.UUID) error
}

// SessionRepository manages user login sessions.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *UserSession) error
	GetSession(ctx context.Context, tokenID string) (*UserSession, error)
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error)
	DeleteSession(ctx context.Context, tokenID string) error
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
	UpdateSessionLastUsed(ctx context.Context, tokenID string) error
}

// PaymentMethodRepository manages the user's stored payment methods.
type PaymentMethodRepository interface {
	CreateUserPaymentMethod(ctx context.Context, pm *UserPaymentMethod) error
	GetUserPaymentMethod(ctx context.Context, pmID uuid.UUID) (*UserPaymentMethod, error)
	GetUserPaymentMethods(ctx context.Context, userID uuid.UUID) ([]UserPaymentMethod, error)
	UpdateUserPaymentMethod(ctx context.Context, pm *UserPaymentMethod) error
	DeleteUserPaymentMethod(ctx context.Context, pmID uuid.UUID) error
	SetDefaultUserPaymentMethod(ctx context.Context, userID, pmID uuid.UUID) error
}

// ─── Composite interface ──────────────────────────────────────────────────────
// UserRepository is the full repository contract required by authService.
// Infrastructure implementations (e.g. the Postgres repository) must satisfy
// this composite interface. All sub-interfaces above are embedded here so that
// existing code that depends on UserRepository continues to compile unchanged.
type UserRepository interface {
	UserCoreRepository
	OAuthRepository
	EmailVerificationRepository
	PasswordResetRepository
	TwoFactorRepository
	SessionRepository
	PaymentMethodRepository
}
