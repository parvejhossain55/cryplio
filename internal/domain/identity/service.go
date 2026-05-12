package identity

// service.go is the entry-point for the identity domain service.
// Method implementations are split across focused files in this package:
//
//   auth.go        — Register, Login, Logout, RefreshToken, Complete2FALogin
//   oauth.go       — Google OAuth (URL generation + callback)
//   profile.go     — User profile reads and updates
//   email.go       — Email verification and password reset
//   twofactor.go   — TOTP 2FA setup, verification, disable
//   session.go     — Session CRUD
//   payment.go     — User payment methods
//   admin.go       — Admin user management
//   helpers.go     — Shared package-level helpers (hashing, validation)

import (
	"context"
	"time"

	"cryplio/internal/domain/events"
	"cryplio/internal/domain/wallet"

	"github.com/google/uuid"
)

// ─── Token type constants ─────────────────────────────────────────────────────

const (
	TokenTypeAccess     = "access"
	TokenTypeRefresh    = "refresh"
	TokenType2FAPending = "2fa_pending"
)

// ─── Error types ──────────────────────────────────────────────────────────────

// TwoFactorRequiredError is returned by Login when the account has 2FA enabled.
// The caller must redirect the user to the 2FA completion step.
type TwoFactorRequiredError struct {
	User      *User
	TempToken string
}

func (e TwoFactorRequiredError) Error() string { return "two factor authentication required" }

// ─── Narrow domain interfaces ─────────────────────────────────────────────────

// UserRegistrar handles user registration.
type UserRegistrar interface {
	Register(ctx context.Context, email, username, password string) (*User, error)
}

// Authenticator handles login/logout and token rotation.
type Authenticator interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *User, err error)
	Logout(ctx context.Context, tokenID string) error
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, user *User, err error)
	Complete2FALogin(ctx context.Context, tempToken, code string) (accessToken, refreshToken string, user *User, err error)
}

// OAuthProvider handles third-party OAuth authentication.
type OAuthProvider interface {
	GoogleOAuthURL() string
	GoogleCallback(ctx context.Context, code string) (accessToken, refreshToken string, user *User, err error)
}

// ProfileManager handles user profile queries and updates.
type ProfileManager interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, *UserStats, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, username, bio *string) (*User, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*User, error)
	UpdateLastSeen(ctx context.Context, userID uuid.UUID) error
}

// EmailVerifier handles email verification flows.
type EmailVerifier interface {
	RequestEmailVerification(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) (*User, error)
}

// PasswordResetter handles password reset flows.
type PasswordResetter interface {
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) (*User, error)
}

// EmailMailer sends transactional account emails.
type EmailMailer interface {
	SendPasswordReset(ctx context.Context, email, token string) error
	SendVerificationEmail(ctx context.Context, email, token string) error
}

// TwoFactorManager handles TOTP 2FA lifecycle.
type TwoFactorManager interface {
	Setup2FA(ctx context.Context, userID uuid.UUID) (secret, provisioningURI string, err error)
	Verify2FA(ctx context.Context, userID uuid.UUID, code string) error
	Disable2FA(ctx context.Context, userID uuid.UUID, password string) error
	Is2FAEnabled(user *User) bool
}

// SessionManager handles user login session CRUD.
type SessionManager interface {
	CreateSession(ctx context.Context, userID uuid.UUID, tokenID string, deviceFingerprint, ipAddress, userAgent, deviceType, location *string, isRemembered bool, expiresAt time.Time) (*UserSession, error)
	GetSession(ctx context.Context, tokenID string) (*UserSession, error)
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error)
	DeleteSession(ctx context.Context, tokenID string) error
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
}

// PaymentManager manages a user's stored payment method profiles.
type PaymentManager interface {
	AddPaymentMethod(ctx context.Context, userID uuid.UUID, pm *UserPaymentMethod) (*UserPaymentMethod, error)
	GetPaymentMethods(ctx context.Context, userID uuid.UUID) ([]UserPaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, userID uuid.UUID, pm *UserPaymentMethod) (*UserPaymentMethod, error)
	RemovePaymentMethod(ctx context.Context, userID, pmID uuid.UUID) error
	SetDefaultPaymentMethod(ctx context.Context, userID, pmID uuid.UUID) error
}

// UserAdminManager handles admin-level user management.
type UserAdminManager interface {
	ListUsers(ctx context.Context, limit, offset int) ([]User, error)
	CountUsers(ctx context.Context) (int, error)
	SuspendUser(ctx context.Context, adminID, userID uuid.UUID, reason string, duration *time.Duration) error
	UnsuspendUser(ctx context.Context, adminID, userID uuid.UUID) error
	BanUser(ctx context.Context, adminID, userID uuid.UUID, reason string) error
	UnbanUser(ctx context.Context, adminID, userID uuid.UUID) error
}

// ─── Composite interface ───────────────────────────────────────────────────────

// AuthService is the full identity service used by the HTTP layer.
// Prefer injecting the narrower interfaces above where possible.
type AuthService interface {
	UserRegistrar
	Authenticator
	OAuthProvider
	ProfileManager
	EmailVerifier
	PasswordResetter
	TwoFactorManager
	SessionManager
	PaymentManager
	UserAdminManager
}

// ─── Dashboard stats (shared type) ───────────────────────────────────────────

// DashboardStats holds aggregated admin dashboard metrics.
type DashboardStats struct {
	TotalUsers      int `json:"total_users"`
	TotalTrades     int `json:"total_trades"`
	PendingTrades   int `json:"pending_trades"`
	ActiveTrades    int `json:"active_trades"`
	PaidTrades      int `json:"paid_trades"`
	CompletedTrades int `json:"completed_trades"`
	DisputedTrades  int `json:"disputed_trades"`
	CancelledTrades int `json:"cancelled_trades"`
	TotalDisputes   int `json:"total_disputes"`
	PendingDisputes int `json:"pending_disputes"`
}

// ─── Concrete service struct ──────────────────────────────────────────────────

type authService struct {
	userRepo           UserRepository
	walletService      wallet.Service
	jwtSecret          string
	jwtExpiry          time.Duration
	refreshTokenExpiry time.Duration
	cookieName         string
	cookieSecure       bool
	cookieSameSite     string
	googleClientID     string
	googleClientSecret string
	oauthRedirectURL   string
	eventDispatcher    events.Dispatcher
	emailMailer        EmailMailer
	issuerName         string // TOTP issuer shown in authenticator apps
}

// NewAuthService constructs an authService with the required dependencies.
// Optional capabilities (Google OAuth, event dispatch, email) are added via
// the With* builder methods.
type AuthServiceConfig struct {
	UserRepo           UserRepository
	JWTSecret          string
	JWTExpiry          time.Duration
	RefreshTokenExpiry time.Duration
	CookieName         string
	CookieSecure       bool
	CookieSameSite     string
	IssuerName         string
}

func NewAuthService(cfg AuthServiceConfig) *authService {
	return &authService{
		userRepo:           cfg.UserRepo,
		jwtSecret:          cfg.JWTSecret,
		jwtExpiry:          cfg.JWTExpiry,
		refreshTokenExpiry: cfg.RefreshTokenExpiry,
		cookieName:         cfg.CookieName,
		cookieSecure:       cfg.CookieSecure,
		cookieSameSite:     cfg.CookieSameSite,
		issuerName:         cfg.IssuerName,
	}
}

// WithGoogleOAuth configures Google OAuth 2.0 credentials.
func (s *authService) WithGoogleOAuth(clientID, clientSecret, redirectURL string) *authService {
	s.googleClientID = clientID
	s.googleClientSecret = clientSecret
	s.oauthRedirectURL = redirectURL
	return s
}

// WithEventDispatcher attaches a domain-event dispatcher.
func (s *authService) WithEventDispatcher(dispatcher events.Dispatcher) *authService {
	s.eventDispatcher = dispatcher
	return s
}

// WithPasswordResetMailer attaches the email client used for password-reset and
// verification emails.
func (s *authService) WithPasswordResetMailer(mailer EmailMailer) *authService {
	s.emailMailer = mailer
	return s
}

// WithWalletService attaches the wallet service for auto-creating wallets on registration.
func (s *authService) WithWalletService(ws wallet.Service) *authService {
	s.walletService = ws
	return s
}
