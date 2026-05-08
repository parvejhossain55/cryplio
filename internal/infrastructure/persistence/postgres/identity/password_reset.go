package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// ─── Password Reset Tokens ────────────────────────────────────────────────────

// CreatePasswordResetToken inserts a new token and populates ID + CreatedAt.
func (r *userRepository) CreatePasswordResetToken(ctx context.Context, t *PasswordResetToken) error {
	var ipAddress sql.NullString
	if t.IPAddress != nil {
		ipAddress = sql.NullString{String: *t.IPAddress, Valid: true}
	}
	return r.db.QueryRowContext(ctx, `
		INSERT INTO password_reset_tokens (user_id, token_hash, ip_address, expires_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at`,
		t.UserID, t.TokenHash, ipAddress, t.ExpiresAt,
	).Scan(&t.ID, &t.CreatedAt)
}

// GetPasswordResetToken looks up an unused token by its SHA-256 hash.
func (r *userRepository) GetPasswordResetToken(ctx context.Context, hash string) (*PasswordResetToken, error) {
	var t PasswordResetToken
	var ipAddress sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, ip_address, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1`, hash,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &ipAddress, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query password reset token: %w", err)
	}
	if ipAddress.Valid {
		t.IPAddress = &ipAddress.String
	}
	return &t, nil
}

// GetPasswordResetTokenByUserID returns the most recent unused token for the
// given user, or nil if none exists.
func (r *userRepository) GetPasswordResetTokenByUserID(ctx context.Context, userID uuid.UUID) (*PasswordResetToken, error) {
	var t PasswordResetToken
	var ipAddress sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, ip_address, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE user_id = $1 AND used_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`, userID,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &ipAddress, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query password reset token by user: %w", err)
	}
	if ipAddress.Valid {
		t.IPAddress = &ipAddress.String
	}
	return &t, nil
}

// MarkPasswordResetTokenUsed stamps used_at = NOW() on the token.
func (r *userRepository) MarkPasswordResetTokenUsed(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE id = $1 AND used_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("mark password reset token used: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("token not found or already used")
	}
	return nil
}
