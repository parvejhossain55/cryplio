package identity

import (
	"context"
	"time"

	"cryplio/pkg/apperrors"

	"github.com/google/uuid"
)

// ─── Admin User Management ────────────────────────────────────────────────────

// CountUsers returns the total number of registered users.
func (s *authService) CountUsers(ctx context.Context) (int, error) {
	return s.userRepo.CountUsers(ctx)
}

// ListUsers returns a paginated slice of users. Limit is clamped to [1, 100].
func (s *authService) ListUsers(ctx context.Context, limit, offset int) ([]User, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return s.userRepo.GetAll(ctx, limit, offset)
}

// SuspendUser suspends a non-admin user account for an optional duration.
// If duration is nil the suspension is indefinite.
func (s *authService) SuspendUser(ctx context.Context, adminID, userID uuid.UUID, reason string, duration *time.Duration) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.Internal("failed to get user", err)
	}
	if user == nil {
		return apperrors.NotFound("user not found", nil)
	}
	if user.Role == UserRoleAdmin {
		return apperrors.Forbidden("cannot suspend admin users", nil)
	}

	var until *time.Time
	if duration != nil {
		t := time.Now().Add(*duration)
		until = &t
	}
	user.Suspend(reason, until)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to suspend user", err)
	}
	return nil
}

// UnsuspendUser lifts an active suspension.
func (s *authService) UnsuspendUser(ctx context.Context, adminID, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.Internal("failed to get user", err)
	}
	if user == nil {
		return apperrors.NotFound("user not found", nil)
	}
	user.Unsuspend()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to unsuspend user", err)
	}
	return nil
}

// BanUser permanently bans a non-admin user and revokes all active sessions.
func (s *authService) BanUser(ctx context.Context, adminID, userID uuid.UUID, reason string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.Internal("failed to get user", err)
	}
	if user == nil {
		return apperrors.NotFound("user not found", nil)
	}
	if user.Role == UserRoleAdmin {
		return apperrors.Forbidden("cannot ban admin users", nil)
	}

	now := time.Now()
	user.Status = UserStatusBanned
	user.IsSuspended = true
	user.SuspensionReason = &reason
	user.SuspendedAt = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to ban user", err)
	}
	// Revoke all sessions so the user is immediately signed out everywhere.
	_ = s.userRepo.DeleteSessionsByUserID(ctx, userID)
	return nil
}

// UnbanUser restores a banned user to active status.
func (s *authService) UnbanUser(ctx context.Context, adminID, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NotFound("user not found", nil)
	}

	user.Status = UserStatusActive
	user.IsSuspended = false
	user.SuspensionReason = nil
	user.SuspendedAt = nil
	user.SuspendedUntil = nil

	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to unban user", err)
	}
	return nil
}
