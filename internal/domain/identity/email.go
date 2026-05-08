package identity

import (
	"context"
	"database/sql"
	"time"

	"cryplio/internal/domain/events"
	"cryplio/pkg/apperrors"
	sharedcrypto "cryplio/pkg/crypto"

	"github.com/google/uuid"
)

// ─── Email Verification ───────────────────────────────────────────────────────

// RequestEmailVerification generates a verification token and sends it to the
// user's email address. Returns an error if the email is already verified.
func (s *authService) RequestEmailVerification(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NotFound("user not found", err)
	}
	if user.EmailVerified {
		return apperrors.Conflict("email already verified", nil)
	}

	plainToken, err := s.generateRandomToken()
	if err != nil {
		return apperrors.Internal("failed to generate token", err)
	}
	tokenHash := hashToken(plainToken)

	token := &EmailVerificationToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.userRepo.CreateEmailVerificationToken(ctx, token); err != nil {
		return apperrors.Internal("failed to create verification token", err)
	}

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.EmailVerificationRequestedEvent{
			UserID:  userID.String(),
			Email:   user.Email,
			Token:   plainToken,
			Expires: token.ExpiresAt,
		})
	}
	if s.emailMailer != nil {
		if err := s.emailMailer.SendVerificationEmail(ctx, user.Email, plainToken); err != nil {
			return apperrors.Internal("failed to send verification email", err)
		}
	}
	return nil
}

// VerifyEmail marks the user's email as verified using the token from the link.
func (s *authService) VerifyEmail(ctx context.Context, tokenString string) (*User, error) {
	tokenHash := hashToken(tokenString)
	token, err := s.userRepo.GetEmailVerificationTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.Internal("database error", err)
	}
	if token == nil {
		return nil, apperrors.Unauthorized("invalid or expired token", nil)
	}
	if token.IsExpired() {
		return nil, apperrors.Unauthorized("token expired", nil)
	}

	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if token.IsVerified() || user.EmailVerified {
		return user, nil // idempotent
	}

	if err := s.userRepo.MarkEmailVerificationTokenVerified(ctx, token.ID); err != nil {
		return nil, apperrors.Internal("failed to mark token verified", err)
	}
	user.EmailVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.EmailVerifiedEvent{
			UserID:     user.UserID.String(),
			Email:      user.Email,
			VerifiedAt: time.Now(),
		})
	}
	return user, nil
}

// ─── Password Reset ───────────────────────────────────────────────────────────

// RequestPasswordReset sends a reset link to the given email if a matching
// account exists. Deliberately returns nil on a missing email (no enumeration).
func (s *authService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // don't reveal whether the email exists
		}
		return apperrors.Internal("database error", err)
	}

	plainToken, err := s.generateRandomToken()
	if err != nil {
		return apperrors.Internal("failed to generate token", err)
	}
	tokenHash := hashToken(plainToken)

	token := &PasswordResetToken{
		UserID:    user.UserID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.userRepo.CreatePasswordResetToken(ctx, token); err != nil {
		return apperrors.Internal("failed to create reset token", err)
	}

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.PasswordResetRequestedEvent{
			UserID: user.UserID.String(),
			Email:  user.Email,
			Token:  plainToken,
		})
	}
	if s.emailMailer != nil {
		if err := s.emailMailer.SendPasswordReset(ctx, user.Email, plainToken); err != nil {
			return apperrors.Internal("failed to send password reset email", err)
		}
	}
	return nil
}

// ResetPassword sets a new password using the single-use token from the link.
// It also invalidates all active sessions as a security measure.
func (s *authService) ResetPassword(ctx context.Context, tokenString, newPassword string) (*User, error) {
	if err := validatePasswordComplexity(newPassword); err != nil {
		return nil, err
	}

	tokenHash := hashToken(tokenString)
	token, err := s.userRepo.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.Internal("database error", err)
	}
	if token == nil {
		return nil, apperrors.Unauthorized("invalid or expired token", nil)
	}
	if token.IsExpired() {
		return nil, apperrors.Unauthorized("token expired", nil)
	}
	if token.IsUsed() {
		return nil, apperrors.Conflict("token already used", nil)
	}

	hashedPassword, err := sharedcrypto.HashPassword(newPassword)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}

	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	user.PasswordHash = hashedPassword
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	if err := s.userRepo.MarkPasswordResetTokenUsed(ctx, token.ID); err != nil {
		return nil, apperrors.Internal("failed to mark token used", err)
	}

	// Invalidate all sessions so stolen refresh tokens become worthless.
	_ = s.userRepo.DeleteSessionsByUserID(ctx, user.UserID)

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.PasswordResetCompletedEvent{
			UserID:  user.UserID.String(),
			Email:   user.Email,
			ResetAt: time.Now(),
		})
	}
	return user, nil
}
