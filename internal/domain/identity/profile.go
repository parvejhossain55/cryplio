package identity

import (
	"context"

	"cryplio/pkg/apperrors"

	"github.com/google/uuid"
)

// ─── Profile queries ──────────────────────────────────────────────────────────

// GetUserByID returns a non-deleted user by their UUID.
func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return nil, apperrors.NotFound("user not found", nil)
	}
	return user, nil
}

// GetUserByUsername returns a user and their trade stats by username.
func (s *authService) GetUserByUsername(ctx context.Context, username string) (*User, *UserStats, error) {
	user, stats, err := s.userRepo.GetByUsernameWithStats(ctx, username)
	if err != nil {
		return nil, nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return nil, nil, apperrors.NotFound("user not found", nil)
	}
	return user, stats, nil
}

// GetUserStats returns the trade statistics for a user.
func (s *authService) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	return s.userRepo.GetStats(ctx, userID)
}

// ─── Profile mutations ────────────────────────────────────────────────────────

// UpdateProfile sets new values for username and/or bio. Nil pointers leave
// the existing value unchanged.
func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, username, bio *string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return nil, apperrors.NotFound("user not found", nil)
	}

	if username != nil && *username != "" {
		user.Username = *username
	}
	if bio != nil {
		user.Bio = bio
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	return user, nil
}

// UpdateAvatar stores a new avatar URL on the user record.
func (s *authService) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return nil, apperrors.NotFound("user not found", nil)
	}

	user.AvatarURL = &avatarURL
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	return user, nil
}

// UpdateLastSeen refreshes the user's last-seen timestamp (best-effort).
func (s *authService) UpdateLastSeen(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.UpdateLastSeen(ctx, userID)
}
