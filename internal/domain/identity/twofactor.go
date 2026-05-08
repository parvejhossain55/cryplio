package identity

import (
	"context"
	"time"

	"cryplio/internal/domain/events"
	"cryplio/pkg/apperrors"
	sharedcrypto "cryplio/pkg/crypto"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

// totpValidate is a thin wrapper so auth.go can validate TOTP codes without
// importing the otp package directly (keeping that dependency in one place).
func totpValidate(code, secret string) bool {
	return totp.Validate(code, secret)
}

// ─── 2FA Lifecycle ────────────────────────────────────────────────────────────

// Is2FAEnabled reports whether the user has an active 2FA secret.
func (s *authService) Is2FAEnabled(user *User) bool {
	return user.TwoFASecret != nil && *user.TwoFASecret != ""
}

// Setup2FA generates a new TOTP secret and provisioning URI. The secret is
// stored as a pending record (10-minute TTL) until confirmed with Verify2FA.
func (s *authService) Setup2FA(ctx context.Context, userID uuid.UUID) (secret, provisioningURI string, err error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", apperrors.NotFound("user not found", err)
	}
	if s.Is2FAEnabled(user) {
		return "", "", apperrors.Conflict("2FA already enabled", nil)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuerName,
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", apperrors.Internal("failed to generate 2FA secret", err)
	}

	pending := &TwoFactorPending{
		ID:        uuid.New(),
		UserID:    userID,
		Secret:    key.Secret(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.userRepo.CreateTwoFactorPending(ctx, pending); err != nil {
		return "", "", apperrors.Internal("failed to store pending 2FA", err)
	}

	return key.Secret(), key.URL(), nil
}

// Verify2FA validates the TOTP code against the pending secret and, on
// success, enables 2FA permanently on the account.
func (s *authService) Verify2FA(ctx context.Context, userID uuid.UUID, code string) error {
	pending, err := s.userRepo.GetTwoFactorPendingByUserID(ctx, userID)
	if err != nil {
		return apperrors.Internal("database error", err)
	}
	if pending == nil {
		return apperrors.Unauthorized("2FA setup not initiated or expired", nil)
	}
	if time.Now().After(pending.ExpiresAt) {
		_ = s.userRepo.DeleteTwoFactorPending(ctx, userID)
		return apperrors.Unauthorized("2FA setup expired, please start again", nil)
	}
	if !totp.Validate(code, pending.Secret) {
		return apperrors.Unauthorized("invalid 2FA code", nil)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NotFound("user not found", err)
	}
	user.TwoFASecret = &pending.Secret
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to update user", err)
	}
	_ = s.userRepo.DeleteTwoFactorPending(ctx, userID)

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.TwoFactorEnabledEvent{
			UserID:    user.UserID.String(),
			EnabledAt: time.Now(),
		})
	}
	return nil
}

// Disable2FA requires the user's password for confirmation, then removes the
// stored 2FA secret.
func (s *authService) Disable2FA(ctx context.Context, userID uuid.UUID, password string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NotFound("user not found", err)
	}
	if err := sharedcrypto.CheckPassword(user.PasswordHash, password); err != nil {
		return apperrors.Unauthorized("invalid password", nil)
	}

	user.TwoFASecret = nil
	_ = s.userRepo.DeleteTwoFactorPending(ctx, userID)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to update user", err)
	}
	return nil
}
