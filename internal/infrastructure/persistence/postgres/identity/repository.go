package identity

// repository.go is the entry point for the identity Postgres repository.
// The UserRepository implementation is split across focused files:
//
//   user.go              — User CRUD + stats + login tracking
//   oauth.go             — OAuth provider link CRUD
//   email_verification.go — Email-verification token CRUD
//   password_reset.go    — Password-reset token CRUD
//   session.go           — Login session CRUD
//   twofactor.go         — In-progress 2FA setup (pending) CRUD
//   payment_method.go    — User payment method profile CRUD
//   helpers.go           — Shared scan utilities

import (
	"database/sql"

	domainidentity "cryplio/internal/domain/identity"
)

// ─── Domain type aliases ──────────────────────────────────────────────────────
// Bring domain types into this package's namespace so every sub-file can use
// them without repeating the import alias.

type User = domainidentity.User
type UserStats = domainidentity.UserStats
type UserOAuth = domainidentity.UserOAuth
type NullUUID = domainidentity.NullUUID
type EmailVerificationToken = domainidentity.EmailVerificationToken
type PasswordResetToken = domainidentity.PasswordResetToken
type UserSession = domainidentity.UserSession
type TwoFactorPending = domainidentity.TwoFactorPending
type UserPaymentMethod = domainidentity.UserPaymentMethod

// UserRepository is the domain interface (aliased for convenience).
type UserRepository = domainidentity.UserRepository

// ─── Concrete struct ──────────────────────────────────────────────────────────

// userRepository implements identity.UserRepository on top of PostgreSQL.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a userRepository backed by the given *sql.DB.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
