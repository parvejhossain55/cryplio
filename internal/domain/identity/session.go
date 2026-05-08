package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ─── Session CRUD ─────────────────────────────────────────────────────────────

// CreateSession stores a new login session with optional device metadata.
func (s *authService) CreateSession(
	ctx context.Context,
	userID uuid.UUID,
	tokenID string,
	deviceFingerprint, ipAddress, userAgent, deviceType, location *string,
	isRemembered bool,
	expiresAt time.Time,
) (*UserSession, error) {
	now := time.Now()
	session := &UserSession{
		ID:                uuid.New(),
		UserID:            userID,
		TokenID:           tokenID,
		DeviceFingerprint: deviceFingerprint,
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		DeviceType:        deviceType,
		Location:          location,
		IsRemembered:      isRemembered,
		ExpiresAt:         expiresAt,
		LastUsedAt:        now,
		CreatedAt:         now,
	}
	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// GetSession retrieves a session by its token ID.
func (s *authService) GetSession(ctx context.Context, tokenID string) (*UserSession, error) {
	return s.userRepo.GetSession(ctx, tokenID)
}

// GetSessionsByUserID returns all active sessions for a user.
func (s *authService) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error) {
	return s.userRepo.GetSessionsByUserID(ctx, userID)
}

// DeleteSession revokes a single session (forces sign-out on that device).
func (s *authService) DeleteSession(ctx context.Context, tokenID string) error {
	return s.userRepo.DeleteSession(ctx, tokenID)
}

// DeleteSessionsByUserID revokes all sessions for a user (e.g. after a ban or
// password reset).
func (s *authService) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteSessionsByUserID(ctx, userID)
}
