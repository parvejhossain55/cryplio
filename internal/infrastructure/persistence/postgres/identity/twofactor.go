package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// ─── Two-Factor Pending ───────────────────────────────────────────────────────
//
// A TwoFactorPending row holds the unconfirmed TOTP secret while the user scans
// the QR code. It is cleaned up once Verify2FA succeeds or the TTL expires.

// CreateTwoFactorPending upserts the pending record (one row per user, replaced
// on re-setup).
func (r *userRepository) CreateTwoFactorPending(ctx context.Context, p *TwoFactorPending) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO two_factor_pending (id, user_id, secret, created_at, expires_at)
		VALUES ($1, $2, $3, NOW(), $4)
		ON CONFLICT (user_id) DO UPDATE
		SET secret     = EXCLUDED.secret,
		    created_at = NOW(),
		    expires_at = EXCLUDED.expires_at`,
		p.ID, p.UserID, p.Secret, p.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("create two factor pending: %w", err)
	}
	return nil
}

// GetTwoFactorPendingByUserID returns the pending record for the user, or nil.
func (r *userRepository) GetTwoFactorPendingByUserID(ctx context.Context, userID uuid.UUID) (*TwoFactorPending, error) {
	var p TwoFactorPending
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, secret, created_at, expires_at
		FROM two_factor_pending
		WHERE user_id = $1`, userID,
	).Scan(&p.ID, &p.UserID, &p.Secret, &p.CreatedAt, &p.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query two factor pending: %w", err)
	}
	return &p, nil
}

// DeleteTwoFactorPending removes the pending record (called after success or expiry).
func (r *userRepository) DeleteTwoFactorPending(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM two_factor_pending WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("delete two factor pending: %w", err)
	}
	return nil
}
