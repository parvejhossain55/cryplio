package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// ─── OAuth Provider Links ─────────────────────────────────────────────────────

// GetOAuthByProviderID returns the OAuth link for a given provider + provider user ID pair.
func (r *userRepository) GetOAuthByProviderID(ctx context.Context, provider, providerUserID string) (*UserOAuth, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_email, provider_username,
		       access_token, refresh_token, token_expiry, created_at, updated_at
		FROM user_oauth
		WHERE provider = $1 AND provider_user_id = $2
	`
	var o UserOAuth
	err := r.db.QueryRowContext(ctx, query, provider, providerUserID).Scan(
		&o.ID, &o.UserID, &o.Provider, &o.ProviderUserID,
		&o.ProviderEmail, &o.ProviderUsername,
		&o.AccessToken, &o.RefreshToken, &o.TokenExpiry,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query oauth by provider: %w", err)
	}
	return &o, nil
}

// CreateOAuth upserts an OAuth link. On conflict it refreshes the stored tokens.
func (r *userRepository) CreateOAuth(ctx context.Context, o *UserOAuth) error {
	query := `
		INSERT INTO user_oauth (
			id, user_id, provider, provider_user_id,
			provider_email, provider_username,
			access_token, refresh_token, token_expiry,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (provider, provider_user_id) DO UPDATE SET
			user_id           = EXCLUDED.user_id,
			provider_email    = EXCLUDED.provider_email,
			provider_username = EXCLUDED.provider_username,
			access_token      = EXCLUDED.access_token,
			refresh_token     = EXCLUDED.refresh_token,
			token_expiry      = EXCLUDED.token_expiry,
			updated_at        = NOW()
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		o.ID, o.UserID, o.Provider, o.ProviderUserID,
		o.ProviderEmail, o.ProviderUsername,
		o.AccessToken, o.RefreshToken, o.TokenExpiry,
	).Scan(&o.ID)
}

// UpdateOAuth refreshes the stored access/refresh tokens and expiry.
func (r *userRepository) UpdateOAuth(ctx context.Context, o *UserOAuth) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_oauth
		SET access_token = $1, refresh_token = $2, token_expiry = $3, updated_at = NOW()
		WHERE id = $4`,
		o.AccessToken, o.RefreshToken, o.TokenExpiry, o.ID,
	)
	if err != nil {
		return fmt.Errorf("update oauth: %w", err)
	}
	return nil
}

// GetOAuthByUserID returns all OAuth links for a user, most recent first.
func (r *userRepository) GetOAuthByUserID(ctx context.Context, userID uuid.UUID) ([]UserOAuth, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, provider, provider_user_id, provider_email, provider_username,
		       access_token, refresh_token, token_expiry, created_at, updated_at
		FROM user_oauth
		WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("query oauth by user: %w", err)
	}
	defer rows.Close()

	var links []UserOAuth
	for rows.Next() {
		var o UserOAuth
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Provider, &o.ProviderUserID,
			&o.ProviderEmail, &o.ProviderUsername,
			&o.AccessToken, &o.RefreshToken, &o.TokenExpiry,
			&o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan oauth: %w", err)
		}
		links = append(links, o)
	}
	return links, nil
}

// DeleteOAuth removes a single OAuth link by its UUID.
func (r *userRepository) DeleteOAuth(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_oauth WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete oauth: %w", err)
	}
	return nil
}
