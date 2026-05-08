package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// ─── Email Verification Tokens ────────────────────────────────────────────────

// CreateEmailVerificationToken inserts a new token and populates ID + CreatedAt.
func (r *userRepository) CreateEmailVerificationToken(ctx context.Context, t *EmailVerificationToken) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO email_verification_tokens (user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at`,
		t.UserID, t.TokenHash, t.ExpiresAt,
	).Scan(&t.ID, &t.CreatedAt)
}

// GetEmailVerificationTokenByHash looks up a token by its SHA-256 hash.
func (r *userRepository) GetEmailVerificationTokenByHash(ctx context.Context, hash string) (*EmailVerificationToken, error) {
	var t EmailVerificationToken
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, verified_at, created_at
		FROM email_verification_tokens
		WHERE token_hash = $1`, hash,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.VerifiedAt, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query email verification token by hash: %w", err)
	}
	return &t, nil
}

// GetEmailVerificationToken looks up a token by its integer primary key.
// Kept for internal use; prefer GetEmailVerificationTokenByHash externally.
func (r *userRepository) GetEmailVerificationToken(ctx context.Context, id int) (*EmailVerificationToken, error) {
	var t EmailVerificationToken
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, verified_at, created_at
		FROM email_verification_tokens
		WHERE id = $1`, id,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.VerifiedAt, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query email verification token by id: %w", err)
	}
	return &t, nil
}

// GetEmailVerificationTokenByUserID returns the most recent unverified token
// for the given user, or nil if none exists.
func (r *userRepository) GetEmailVerificationTokenByUserID(ctx context.Context, userID uuid.UUID) (*EmailVerificationToken, error) {
	var t EmailVerificationToken
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, verified_at, created_at
		FROM email_verification_tokens
		WHERE user_id = $1 AND verified_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`, userID,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.VerifiedAt, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query email verification token by user: %w", err)
	}
	return &t, nil
}

// MarkEmailVerificationTokenVerified stamps verified_at = NOW() on the token.
func (r *userRepository) MarkEmailVerificationTokenVerified(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE email_verification_tokens
		SET verified_at = NOW()
		WHERE id = $1 AND verified_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("mark email verification token verified: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("token not found or already verified")
	}
	return nil
}
