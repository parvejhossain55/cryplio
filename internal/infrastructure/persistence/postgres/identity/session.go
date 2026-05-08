package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// ─── Login Sessions ───────────────────────────────────────────────────────────

// CreateSession inserts a new login session row.
func (r *userRepository) CreateSession(ctx context.Context, s *UserSession) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_sessions (
			id, user_id, token_id,
			device_fingerprint, ip_address, user_agent, device_type, location,
			is_remembered, expires_at, last_used_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())`,
		s.ID, s.UserID, s.TokenID,
		s.DeviceFingerprint, s.IPAddress, s.UserAgent, s.DeviceType, s.Location,
		s.IsRemembered, s.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

// GetSession returns a session by its token ID, or nil if not found.
func (r *userRepository) GetSession(ctx context.Context, tokenID string) (*UserSession, error) {
	var s UserSession
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_id,
		       device_fingerprint, ip_address, user_agent, device_type, location,
		       is_remembered, expires_at, last_used_at, created_at
		FROM user_sessions
		WHERE token_id = $1`, tokenID,
	).Scan(
		&s.ID, &s.UserID, &s.TokenID,
		&s.DeviceFingerprint, &s.IPAddress, &s.UserAgent, &s.DeviceType, &s.Location,
		&s.IsRemembered, &s.ExpiresAt, &s.LastUsedAt, &s.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query session: %w", err)
	}
	return &s, nil
}

// GetSessionsByUserID returns all sessions for the user ordered by most-recent-use.
func (r *userRepository) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, token_id,
		       device_fingerprint, ip_address, user_agent, device_type, location,
		       is_remembered, expires_at, last_used_at, created_at
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY last_used_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("query sessions by user: %w", err)
	}
	defer rows.Close()

	var sessions []UserSession
	for rows.Next() {
		var s UserSession
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.TokenID,
			&s.DeviceFingerprint, &s.IPAddress, &s.UserAgent, &s.DeviceType, &s.Location,
			&s.IsRemembered, &s.ExpiresAt, &s.LastUsedAt, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// DeleteSession removes a session by token ID (forces sign-out on that device).
func (r *userRepository) DeleteSession(ctx context.Context, tokenID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_sessions WHERE token_id = $1`, tokenID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// DeleteSessionsByUserID removes all sessions for a user (e.g. after a ban).
func (r *userRepository) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_sessions WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("delete sessions by user: %w", err)
	}
	return nil
}

// UpdateSessionLastUsed stamps last_used_at = NOW() for session activity tracking.
func (r *userRepository) UpdateSessionLastUsed(ctx context.Context, tokenID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_sessions SET last_used_at = NOW() WHERE token_id = $1`, tokenID)
	if err != nil {
		return fmt.Errorf("update session last used: %w", err)
	}
	return nil
}
