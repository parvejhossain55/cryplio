package identity

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"cryplio/internal/domain/events"
	"cryplio/pkg/apperrors"
	sharedcrypto "cryplio/pkg/crypto"
	sharedjwt "cryplio/pkg/jwt"

	"github.com/google/uuid"
)

// ─── Registration ─────────────────────────────────────────────────────────────

// Register creates a new user account after validating the password and
// checking for duplicates.
func (s *authService) Register(ctx context.Context, email, username, password string) (*User, error) {
	if err := validatePasswordComplexity(password); err != nil {
		return nil, err
	}

	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, apperrors.Internal("database error", err)
	}
	if existing != nil {
		return nil, apperrors.Conflict("user already exists", nil)
	}

	hashedPassword, err := sharedcrypto.HashPassword(password)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}

	user := NewUser(email, username, hashedPassword)
	if err := s.userRepo.Create(ctx, user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, apperrors.Conflict("user already exists", nil)
		}
		return nil, apperrors.Internal("failed to create user", err)
	}

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.UserRegisteredEvent{
			UserID:    user.UserID.String(),
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		})
	}

	// Auto-create a default wallet for the new user
	if s.walletService != nil {
		_, _ = s.walletService.CreateDefaultWallet(ctx, user.UserID)
	}

	return user, nil
}

// ─── Login ────────────────────────────────────────────────────────────────────

// Login authenticates a user. If 2FA is enabled it returns a
// TwoFactorRequiredError containing a short-lived temp token; otherwise it
// completes the login and returns access + refresh tokens.
func (s *authService) Login(ctx context.Context, email, password string) (string, string, *User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil, apperrors.Unauthorized("invalid credentials", nil)
		}
		return "", "", nil, apperrors.Internal("database error", err)
	}
	if user == nil {
		return "", "", nil, apperrors.Unauthorized("invalid credentials", nil)
	}

	if user.IsLocked() {
		return "", "", nil, apperrors.RateLimited("account temporarily locked", nil)
	}

	if err := sharedcrypto.CheckPassword(user.PasswordHash, password); err != nil {
		_, _ = s.userRepo.IncrementFailedAttempts(ctx, user.UserID)
		return "", "", nil, apperrors.Unauthorized("invalid credentials", nil)
	}

	if s.Is2FAEnabled(user) {
		tempToken, err := s.generate2FATempToken(user)
		if err != nil {
			return "", "", nil, apperrors.Internal("failed to generate 2FA token", err)
		}
		return "", "", user, TwoFactorRequiredError{User: user, TempToken: tempToken}
	}

	access, refresh, err := s.completeLogin(ctx, user)
	if err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

// ─── Logout ───────────────────────────────────────────────────────────────────

// Logout deletes the user's session and emits a UserLoggedOutEvent.
func (s *authService) Logout(ctx context.Context, tokenID string) error {
	if tokenID != "" {
		_ = s.userRepo.DeleteSession(ctx, tokenID)
	}
	if s.eventDispatcher != nil {
		if userID, ok := ctx.Value("user_id").(string); ok {
			_ = s.eventDispatcher.Dispatch(ctx, events.UserLoggedOutEvent{
				UserID:      userID,
				LoggedOutAt: time.Now(),
			})
		}
	}
	return nil
}

// ─── Token Refresh ────────────────────────────────────────────────────────────

// RefreshToken validates a refresh token, rotates the session, and issues new
// access + refresh tokens.
func (s *authService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, *User, error) {
	claims, err := sharedjwt.Parse(s.jwtSecret, refreshTokenString)
	if err != nil {
		return "", "", nil, apperrors.Unauthorized("invalid refresh token", err)
	}
	tokenType, ok := claims[sharedjwt.ClaimTokenType].(string)
	if !ok || tokenType != TokenTypeRefresh {
		return "", "", nil, apperrors.Unauthorized("invalid token type", nil)
	}
	userIDStr, ok := claims[sharedjwt.ClaimUserID].(string)
	if !ok {
		return "", "", nil, apperrors.Unauthorized("invalid token claims", nil)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", nil, apperrors.InvalidInput("invalid user id", err)
	}
	jti, ok := claims["jti"].(string)
	if !ok {
		return "", "", nil, apperrors.Unauthorized("invalid token", nil)
	}

	session, err := s.userRepo.GetSession(ctx, jti)
	if err != nil {
		return "", "", nil, apperrors.Internal("database error", err)
	}
	if session == nil {
		return "", "", nil, apperrors.Unauthorized("invalid session", nil)
	}
	if session.IsExpired() {
		_ = s.userRepo.DeleteSession(ctx, jti)
		return "", "", nil, apperrors.Unauthorized("session expired", nil)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", nil, apperrors.Unauthorized("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return "", "", nil, apperrors.Unauthorized("user not found", nil)
	}

	// Rotate: delete old session, create new one.
	_ = s.userRepo.DeleteSession(ctx, jti)
	newSessionID := uuid.New()
	newSession := &UserSession{
		ID:           newSessionID,
		UserID:       user.UserID,
		TokenID:      newSessionID.String(),
		IsRemembered: session.IsRemembered,
		ExpiresAt:    time.Now().Add(s.refreshTokenExpiry),
	}
	if err := s.userRepo.CreateSession(ctx, newSession); err != nil {
		return "", "", nil, apperrors.Internal("failed to create session", err)
	}

	accessToken, err := s.generateJWT(user)
	if err != nil {
		return "", "", nil, apperrors.Internal("failed to generate access token", err)
	}
	refreshToken, err := s.generateRefreshJWT(user, newSessionID.String())
	if err != nil {
		return "", "", nil, apperrors.Internal("failed to generate refresh token", err)
	}
	return accessToken, refreshToken, user, nil
}

// ─── 2FA Login Completion ─────────────────────────────────────────────────────

// Complete2FALogin validates the temp token and TOTP code, then finishes the
// login process by issuing access + refresh tokens.
func (s *authService) Complete2FALogin(ctx context.Context, tempToken, code string) (string, string, *User, error) {
	claims, err := sharedjwt.Parse(s.jwtSecret, tempToken)
	if err != nil {
		return "", "", nil, apperrors.Unauthorized("invalid verification token", err)
	}
	tokenType, ok := claims[sharedjwt.ClaimTokenType].(string)
	if !ok || tokenType != TokenType2FAPending {
		return "", "", nil, apperrors.Unauthorized("invalid token type", nil)
	}
	userIDStr, ok := claims[sharedjwt.ClaimUserID].(string)
	if !ok {
		return "", "", nil, apperrors.Unauthorized("invalid token claims", nil)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", nil, apperrors.InvalidInput("invalid user id", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", nil, apperrors.NotFound("user not found", err)
	}
	if !s.Is2FAEnabled(user) {
		return "", "", nil, apperrors.Unauthorized("2FA not enabled", nil)
	}
	if code == "" {
		return "", "", nil, apperrors.Validation("verification code required", nil)
	}
	if user.TwoFASecret == nil || *user.TwoFASecret == "" {
		return "", "", nil, apperrors.Unauthorized("2FA secret not configured", nil)
	}

	// Import totp inline to keep this file free of the otp dependency.
	// Validation is delegated to the same TOTP library used by twofactor.go.
	// Because both files are in the same package, we call the package-local
	// helper defined in twofactor.go.
	if !totpValidate(code, *user.TwoFASecret) {
		return "", "", nil, apperrors.Unauthorized("invalid verification code", nil)
	}

	access, refresh, err := s.completeLogin(ctx, user)
	if err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// completeLogin is called after all authentication checks pass.
// It creates a session, generates tokens, updates login stats, and
// publishes a UserLoggedInEvent.
func (s *authService) completeLogin(ctx context.Context, user *User) (string, string, error) {
	_ = s.userRepo.IncrementLogin(ctx, user.UserID)
	_ = s.userRepo.UpdateLastSeen(ctx, user.UserID)

	sessionID := uuid.New()
	session := &UserSession{
		ID:           sessionID,
		UserID:       user.UserID,
		TokenID:      sessionID.String(),
		IsRemembered: false,
		ExpiresAt:    time.Now().Add(s.refreshTokenExpiry),
	}
	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return "", "", apperrors.Internal("failed to create session", err)
	}

	accessToken, err := s.generateJWT(user)
	if err != nil {
		return "", "", apperrors.Internal("failed to generate access token", err)
	}
	refreshToken, err := s.generateRefreshJWT(user, session.ID.String())
	if err != nil {
		return "", "", apperrors.Internal("failed to generate refresh token", err)
	}

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.UserLoggedInEvent{
			UserID:     user.UserID.String(),
			Email:      user.Email,
			LoggedInAt: time.Now(),
		})
	}

	return accessToken, refreshToken, nil
}

// generateJWT issues a short-lived access token.
func (s *authService) generateJWT(user *User) (string, error) {
	return sharedjwt.Issue(s.jwtSecret, s.jwtExpiry, sharedjwt.Claims{
		sharedjwt.ClaimUserID:    user.UserID.String(),
		"email":                  user.Email,
		"username":               user.Username,
		"role":                   string(user.Role),
		sharedjwt.ClaimTokenType: TokenTypeAccess,
	})
}

// generateRefreshJWT issues a long-lived refresh token bound to a session ID.
func (s *authService) generateRefreshJWT(user *User, sessionID string) (string, error) {
	return sharedjwt.Issue(s.jwtSecret, s.refreshTokenExpiry, sharedjwt.Claims{
		sharedjwt.ClaimUserID:    user.UserID.String(),
		sharedjwt.ClaimTokenType: TokenTypeRefresh,
		"jti":                    sessionID,
	})
}

// generate2FATempToken issues a 5-minute token used during the 2FA challenge.
func (s *authService) generate2FATempToken(user *User) (string, error) {
	return sharedjwt.Issue(s.jwtSecret, 5*time.Minute, sharedjwt.Claims{
		sharedjwt.ClaimUserID:    user.UserID.String(),
		sharedjwt.ClaimTokenType: TokenType2FAPending,
	})
}

// generateRandomToken returns a cryptographically secure 128-bit hex token.
func (s *authService) generateRandomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// generateUniqueUsername tries the base name, then numeric suffixes, then a
// random hex fallback — making at most 9 DB queries.
func (s *authService) generateUniqueUsername(ctx context.Context, base string) string {
	username := strings.ToLower(strings.ReplaceAll(base, " ", "."))

	if existing, _ := s.userRepo.GetByUsername(ctx, username); existing == nil {
		return username
	}
	for i := 2; i <= 9; i++ {
		candidate := fmt.Sprintf("%s%d", username, i)
		if existing, _ := s.userRepo.GetByUsername(ctx, candidate); existing == nil {
			return candidate
		}
	}
	b := make([]byte, 3)
	if _, err := rand.Read(b); err == nil {
		return fmt.Sprintf("%s_%s", username, hex.EncodeToString(b))
	}
	return fmt.Sprintf("%s_%d", username, time.Now().UnixNano()%1000000)
}
