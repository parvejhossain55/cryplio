package identity

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainidentity "cryplio/internal/domain/identity"

	"github.com/google/uuid"
)

// Type aliases for domain types
type User = domainidentity.User
type UserStats = domainidentity.UserStats
type UserOAuth = domainidentity.UserOAuth
type NullUUID = domainidentity.NullUUID
type EmailVerificationToken = domainidentity.EmailVerificationToken
type PasswordResetToken = domainidentity.PasswordResetToken
type UserSession = domainidentity.UserSession
type TwoFactorPending = domainidentity.TwoFactorPending
type UserBlock = domainidentity.UserBlock

// UserRepository is the domain interface (aliased for convenience)
type UserRepository = domainidentity.UserRepository

// userRepository implements UserRepository using PostgreSQL
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// GetAll returns all users (admin use, with pagination)
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]User, error) {
	query := `
		SELECT user_id, email, username, password_hash, phone_country_code, phone_number,
		       phone_verified, email_verified, status,
		       avatar_url, bio, timezone, locale, is_merchant, is_suspended,
		       suspension_reason, suspended_at, suspended_until, last_login_at,
		       login_count, failed_login_attempts, locked_until, referral_code,
		       referred_by, two_fa_secret, remember_2fa, remember_2fa_expiry,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var referredBy sql.NullString
		err := rows.Scan(
			&u.UserID, &u.Email, &u.Username, &u.PasswordHash,
			&u.PhoneCountryCode, &u.PhoneNumber, &u.PhoneVerified, &u.EmailVerified,
			&u.Status, &u.AvatarURL, &u.Bio,
			&u.Timezone, &u.Locale, &u.IsMerchant, &u.IsSuspended,
			&u.SuspensionReason, &u.SuspendedAt, &u.SuspendedUntil, &u.LastLoginAt,
			&u.LoginCount, &u.FailedLoginAttempts, &u.LockedUntil,
			&u.ReferralCode, &referredBy, &u.TwoFASecret, &u.Remember2FA,
			&u.Remember2FAExpiry, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		if referredBy.Valid {
			parsed, err := uuid.Parse(referredBy.String)
			if err == nil {
				u.ReferredBy = NullUUID{UUID: parsed, Valid: true}
			}
		}
		users = append(users, u)
	}
	return users, nil
}

// GetByID returns a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT user_id, email, username, password_hash, phone_country_code, phone_number,
		       phone_verified, email_verified, status,
		       avatar_url, bio, timezone, locale, is_merchant, is_suspended,
		       suspension_reason, suspended_at, suspended_until, last_login_at,
		       login_count, failed_login_attempts, locked_until, referral_code,
		       referred_by, two_fa_secret, remember_2fa, remember_2fa_expiry,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	var u User
	var referredBy sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.UserID, &u.Email, &u.Username, &u.PasswordHash,
		&u.PhoneCountryCode, &u.PhoneNumber, &u.PhoneVerified, &u.EmailVerified,
		&u.Status, &u.AvatarURL, &u.Bio,
		&u.Timezone, &u.Locale, &u.IsMerchant, &u.IsSuspended,
		&u.SuspensionReason, &u.SuspendedAt, &u.SuspendedUntil, &u.LastLoginAt,
		&u.LoginCount, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.ReferralCode, &referredBy, &u.TwoFASecret, &u.Remember2FA,
		&u.Remember2FAExpiry, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	if referredBy.Valid {
		parsed, err := uuid.Parse(referredBy.String)
		if err == nil {
			u.ReferredBy = NullUUID{UUID: parsed, Valid: true}
		}
	}
	return &u, nil
}

// GetByEmail returns a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT user_id, email, username, password_hash, phone_country_code, phone_number,
		       phone_verified, email_verified, status,
		       avatar_url, bio, timezone, locale, is_merchant, is_suspended,
		       suspension_reason, suspended_at, suspended_until, last_login_at,
		       login_count, failed_login_attempts, locked_until, referral_code,
		       referred_by, two_fa_secret, remember_2fa, remember_2fa_expiry,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var u User
	var referredBy sql.NullString
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.UserID, &u.Email, &u.Username, &u.PasswordHash,
		&u.PhoneCountryCode, &u.PhoneNumber, &u.PhoneVerified, &u.EmailVerified,
		&u.Status, &u.AvatarURL, &u.Bio,
		&u.Timezone, &u.Locale, &u.IsMerchant, &u.IsSuspended,
		&u.SuspensionReason, &u.SuspendedAt, &u.SuspendedUntil, &u.LastLoginAt,
		&u.LoginCount, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.ReferralCode, &referredBy, &u.TwoFASecret, &u.Remember2FA,
		&u.Remember2FAExpiry, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by email: %w", err)
	}
	if referredBy.Valid {
		parsed, err := uuid.Parse(referredBy.String)
		if err == nil {
			u.ReferredBy = NullUUID{UUID: parsed, Valid: true}
		}
	}
	return &u, nil
}

// GetByUsername returns a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT user_id, email, username, password_hash, phone_country_code, phone_number,
		       phone_verified, email_verified, status,
		       avatar_url, bio, timezone, locale, is_merchant, is_suspended,
		       suspension_reason, suspended_at, suspended_until, last_login_at,
		       login_count, failed_login_attempts, locked_until, referral_code,
		       referred_by, two_fa_secret, remember_2fa, remember_2fa_expiry,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`
	var u User
	var referredByNull sql.NullString
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.UserID, &u.Email, &u.Username, &u.PasswordHash,
		&u.PhoneCountryCode, &u.PhoneNumber, &u.PhoneVerified, &u.EmailVerified,
		&u.Status, &u.AvatarURL, &u.Bio,
		&u.Timezone, &u.Locale, &u.IsMerchant, &u.IsSuspended,
		&u.SuspensionReason, &u.SuspendedAt, &u.SuspendedUntil, &u.LastLoginAt,
		&u.LoginCount, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.ReferralCode, &referredByNull, &u.TwoFASecret, &u.Remember2FA,
		&u.Remember2FAExpiry, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by username: %w", err)
	}
	if referredByNull.Valid {
		parsed, err := uuid.Parse(referredByNull.String)
		if err == nil {
			u.ReferredBy = NullUUID{UUID: parsed, Valid: true}
		}
	}
	return &u, nil
}

// Create inserts a new user
func (r *userRepository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (
			user_id, email, username, password_hash, phone_country_code, phone_number,
			phone_verified, email_verified, status, timezone, locale,
			is_merchant, login_count, failed_login_attempts, referral_code, remember_2fa,
			referred_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			NOW(), NOW()
		) RETURNING created_at, updated_at
	`
	referredBy := sql.NullString{}
	if u.ReferredBy.Valid {
		referredBy = sql.NullString{String: u.ReferredBy.UUID.String(), Valid: true}
	}
	err := r.db.QueryRowContext(
		ctx, query,
		u.UserID, u.Email, u.Username, u.PasswordHash,
		u.PhoneCountryCode, u.PhoneNumber, u.PhoneVerified, u.EmailVerified,
		u.Status, u.Timezone, u.Locale, u.IsMerchant,
		u.LoginCount, u.FailedLoginAttempts, u.ReferralCode, u.Remember2FA,
		referredBy,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

// Update updates user profile
func (r *userRepository) Update(ctx context.Context, u *User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2, phone_country_code = $3, phone_number = $4,
		    phone_verified = $5, email_verified = $6, avatar_url = $9, bio = $10, timezone = $11, 
		    locale = $12, is_merchant = $13, is_suspended = $14, 
		    suspension_reason = $15, suspended_at = $16, suspended_until = $17, 
		    last_login_at = $18, login_count = $19, failed_login_attempts = $20, 
		    locked_until = $21, referral_code = $22, referred_by = $23, 
		    two_fa_secret = $24, remember_2fa = $25, remember_2fa_expiry = $26, 
		    updated_at = NOW()
		WHERE user_id = $27 AND deleted_at IS NULL
		RETURNING updated_at
	`
	nullReferredBy := sql.NullString{}
	if u.ReferredBy.Valid {
		nullReferredBy = sql.NullString{String: u.ReferredBy.UUID.String(), Valid: true}
	}
	err := r.db.QueryRowContext(
		ctx, query,
		u.Email, u.Username, u.PhoneCountryCode, u.PhoneNumber,
		u.PhoneVerified, u.EmailVerified,
		u.AvatarURL, u.Bio, u.Timezone, u.Locale, u.IsMerchant,
		u.IsSuspended, u.SuspensionReason, u.SuspendedAt, u.SuspendedUntil,
		u.LastLoginAt, u.LoginCount, u.FailedLoginAttempts, u.LockedUntil,
		u.ReferralCode, nullReferredBy, u.TwoFASecret, u.Remember2FA,
		u.Remember2FAExpiry, u.UserID,
	).Scan(&u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete soft-deletes a user
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE user_id = $1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}

// IncrementLogin updates last login info
func (r *userRepository) IncrementLogin(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET last_login_at = NOW(),
		    login_count = login_count + 1,
		    failed_login_attempts = 0,
		    locked_until = NULL,
		    updated_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("update login: %w", err)
	}
	return nil
}

// IncrementFailedAttempts increments failed login counter
func (r *userRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		UPDATE users
		SET failed_login_attempts = failed_login_attempts + 1,
		    locked_until = CASE
		        WHEN failed_login_attempts >= 5 THEN NOW() + INTERVAL '15 minutes'
		        ELSE locked_until
		    END,
		    updated_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL
		RETURNING failed_login_attempts
	`
	var attempts int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&attempts)
	if err != nil {
		return 0, fmt.Errorf("increment failed attempts: %w", err)
	}
	return attempts, nil
}

// GetStats returns user statistics
func (r *userRepository) GetStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	query := `
		SELECT user_id, total_trades, successful_trades, dispute_rate, avg_rating,
		       positive_feedback_count, neutral_feedback_count, negative_feedback_count,
		       total_volume_usd, last_trade_at, updated_at
		FROM user_stats
		WHERE user_id = $1
	`
	var stats UserStats
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.UserID, &stats.TotalTrades, &stats.SuccessfulTrades, &stats.DisputeRate,
		&stats.AvgRating, &stats.PositiveFeedbackCount, &stats.NeutralFeedbackCount,
		&stats.NegativeFeedbackCount, &stats.TotalVolumeUSD, &stats.LastTradeAt, &stats.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user stats: %w", err)
	}
	return &stats, nil
}

// OAuth Methods

// GetOAuthByProviderID returns OAuth account by provider and provider user ID
func (r *userRepository) GetOAuthByProviderID(ctx context.Context, provider, providerUserID string) (*UserOAuth, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_email, provider_username,
		       access_token, refresh_token, token_expiry, created_at, updated_at
		FROM user_oauth
		WHERE provider = $1 AND provider_user_id = $2
	`
	var oauth UserOAuth
	err := r.db.QueryRowContext(ctx, query, provider, providerUserID).Scan(
		&oauth.ID, &oauth.UserID, &oauth.Provider, &oauth.ProviderUserID,
		&oauth.ProviderEmail, &oauth.ProviderUsername,
		&oauth.AccessToken, &oauth.RefreshToken, &oauth.TokenExpiry,
		&oauth.CreatedAt, &oauth.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query oauth: %w", err)
	}
	return &oauth, nil
}

// CreateOAuth creates a new OAuth connection
func (r *userRepository) CreateOAuth(ctx context.Context, oauth *UserOAuth) error {
	query := `
		INSERT INTO user_oauth (id, user_id, provider, provider_user_id, provider_email, provider_username, access_token, refresh_token, token_expiry, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (provider, provider_user_id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			provider_email = EXCLUDED.provider_email,
			provider_username = EXCLUDED.provider_username,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			token_expiry = EXCLUDED.token_expiry,
			updated_at = NOW()
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		oauth.ID, oauth.UserID, oauth.Provider, oauth.ProviderUserID,
		oauth.ProviderEmail, oauth.ProviderUsername,
		oauth.AccessToken, oauth.RefreshToken, oauth.TokenExpiry,
	).Scan(&oauth.ID)
	if err != nil {
		return fmt.Errorf("create oauth: %w", err)
	}
	return nil
}

// UpdateOAuth updates an OAuth account
func (r *userRepository) UpdateOAuth(ctx context.Context, oauth *UserOAuth) error {
	query := `
		UPDATE user_oauth
		SET access_token = $1, refresh_token = $2, token_expiry = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		oauth.AccessToken, oauth.RefreshToken, oauth.TokenExpiry, oauth.ID,
	)
	if err != nil {
		return fmt.Errorf("update oauth: %w", err)
	}
	return nil
}

// GetOAuthByUserID returns all OAuth accounts for a user
func (r *userRepository) GetOAuthByUserID(ctx context.Context, userID uuid.UUID) ([]UserOAuth, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_email, provider_username,
		       access_token, refresh_token, token_expiry, created_at, updated_at
		FROM user_oauth
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query oauth by user: %w", err)
	}
	defer rows.Close()

	var oauthAccounts []UserOAuth
	for rows.Next() {
		var oauth UserOAuth
		err := rows.Scan(
			&oauth.ID, &oauth.UserID, &oauth.Provider, &oauth.ProviderUserID,
			&oauth.ProviderEmail, &oauth.ProviderUsername,
			&oauth.AccessToken, &oauth.RefreshToken, &oauth.TokenExpiry,
			&oauth.CreatedAt, &oauth.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan oauth: %w", err)
		}
		oauthAccounts = append(oauthAccounts, oauth)
	}
	return oauthAccounts, nil
}

// DeleteOAuth removes an OAuth account linkage
func (r *userRepository) DeleteOAuth(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_oauth WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete oauth: %w", err)
	}
	return nil
}

// ============================================================
// Email Verification Token Methods
// ============================================================

// CreateEmailVerificationToken creates a new email verification token
func (r *userRepository) CreateEmailVerificationToken(ctx context.Context, token *EmailVerificationToken) error {
	query := `
		INSERT INTO email_verification_tokens (user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		token.UserID, token.TokenHash, token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
	if err != nil {
		return fmt.Errorf("create email verification token: %w", err)
	}
	return nil
}

// GetEmailVerificationTokenByHash returns a token by its hash
func (r *userRepository) GetEmailVerificationTokenByHash(ctx context.Context, tokenHash string) (*EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, verified_at, created_at
		FROM email_verification_tokens
		WHERE token_hash = $1
	`
	var token EmailVerificationToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt,
		&token.VerifiedAt, &token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query email verification token by hash: %w", err)
	}
	return &token, nil
}

// GetEmailVerificationToken returns a token by ID (optional, kept for compatibility)
func (r *userRepository) GetEmailVerificationToken(ctx context.Context, id int) (*EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, verified_at, created_at
		FROM email_verification_tokens
		WHERE id = $1
	`
	var token EmailVerificationToken
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt,
		&token.VerifiedAt, &token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query email verification token: %w", err)
	}
	return &token, nil
}

// GetEmailVerificationTokenByUserID returns the latest unverified token for a user
func (r *userRepository) GetEmailVerificationTokenByUserID(ctx context.Context, userID uuid.UUID) (*EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, verified_at, created_at
		FROM email_verification_tokens
		WHERE user_id = $1 AND verified_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	var token EmailVerificationToken
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt,
		&token.VerifiedAt, &token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query email verification token by user: %w", err)
	}
	return &token, nil
}

// MarkEmailVerificationTokenVerified marks a token as verified
func (r *userRepository) MarkEmailVerificationTokenVerified(ctx context.Context, id int) error {
	query := `
		UPDATE email_verification_tokens
		SET verified_at = NOW()
		WHERE id = $1 AND verified_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark email verification token verified: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("token not found or already verified")
	}
	return nil
}

// ============================================================
// Password Reset Token Methods
// ============================================================

// CreatePasswordResetToken creates a new password reset token
func (r *userRepository) CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token_hash, ip_address, expires_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	var ipAddress sql.NullString
	if token.IPAddress != nil {
		ipAddress = sql.NullString{String: *token.IPAddress, Valid: true}
	}
	err := r.db.QueryRowContext(
		ctx, query,
		token.UserID, token.TokenHash, ipAddress, token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
	if err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}
	return nil
}

// GetPasswordResetToken returns a token by its hash
func (r *userRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (*PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, ip_address, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1
	`
	var token PasswordResetToken
	var ipAddress sql.NullString
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &ipAddress,
		&token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query password reset token: %w", err)
	}
	if ipAddress.Valid {
		token.IPAddress = &ipAddress.String
	}
	return &token, nil
}

// GetPasswordResetTokenByUserID returns the latest unused token for a user
func (r *userRepository) GetPasswordResetTokenByUserID(ctx context.Context, userID uuid.UUID) (*PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, ip_address, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE user_id = $1 AND used_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	var token PasswordResetToken
	var ipAddress sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &ipAddress,
		&token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query password reset token by user: %w", err)
	}
	if ipAddress.Valid {
		token.IPAddress = &ipAddress.String
	}
	return &token, nil
}

// MarkPasswordResetTokenUsed marks a token as used
func (r *userRepository) MarkPasswordResetTokenUsed(ctx context.Context, id int) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE id = $1 AND used_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark password reset token used: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("token not found or already used")
	}
	return nil
}

// ============================================================
// Session Methods
// ============================================================

// CreateSession creates a new user session
func (r *userRepository) CreateSession(ctx context.Context, session *UserSession) error {
	query := `
		INSERT INTO user_sessions (
			id, user_id, token_id, device_fingerprint, ip_address,
			user_agent, device_type, location, is_remembered, expires_at, last_used_at, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`
	_, err := r.db.ExecContext(
		ctx, query,
		session.ID, session.UserID, session.TokenID, session.DeviceFingerprint,
		session.IPAddress, session.UserAgent, session.DeviceType, session.Location,
		session.IsRemembered, session.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

// GetSession returns a session by token ID
func (r *userRepository) GetSession(ctx context.Context, tokenID string) (*UserSession, error) {
	query := `
		SELECT id, user_id, token_id, device_fingerprint, ip_address, user_agent,
		       device_type, location, is_remembered, expires_at, last_used_at, created_at
		FROM user_sessions
		WHERE token_id = $1
	`
	var session UserSession
	err := r.db.QueryRowContext(ctx, query, tokenID).Scan(
		&session.ID, &session.UserID, &session.TokenID, &session.DeviceFingerprint,
		&session.IPAddress, &session.UserAgent, &session.DeviceType, &session.Location,
		&session.IsRemembered, &session.ExpiresAt, &session.LastUsedAt, &session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query session: %w", err)
	}
	return &session, nil
}

// GetSessionsByUserID returns all sessions for a user
func (r *userRepository) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error) {
	query := `
		SELECT id, user_id, token_id, device_fingerprint, ip_address, user_agent,
		       device_type, location, is_remembered, expires_at, last_used_at, created_at
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY last_used_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query sessions by user: %w", err)
	}
	defer rows.Close()

	var sessions []UserSession
	for rows.Next() {
		var s UserSession
		err := rows.Scan(
			&s.ID, &s.UserID, &s.TokenID, &s.DeviceFingerprint,
			&s.IPAddress, &s.UserAgent, &s.DeviceType, &s.Location,
			&s.IsRemembered, &s.ExpiresAt, &s.LastUsedAt, &s.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// DeleteSession removes a session by token ID
func (r *userRepository) DeleteSession(ctx context.Context, tokenID string) error {
	query := `DELETE FROM user_sessions WHERE token_id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// DeleteSessionsByUserID removes all sessions for a user (except current)
func (r *userRepository) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete sessions by user: %w", err)
	}
	return nil
}

// UpdateSessionLastUsed updates the last_used_at timestamp
func (r *userRepository) UpdateSessionLastUsed(ctx context.Context, tokenID string) error {
	query := `
		UPDATE user_sessions
		SET last_used_at = NOW()
		WHERE token_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("update session last used: %w", err)
	}
	return nil
}

// ============================================================
// Two-Factor Pending Methods
// ============================================================

// CreateTwoFactorPending creates a new pending 2FA record
func (r *userRepository) CreateTwoFactorPending(ctx context.Context, pending *TwoFactorPending) error {
	query := `
		INSERT INTO two_factor_pending (id, user_id, secret, created_at, expires_at)
		VALUES ($1, $2, $3, NOW(), $4)
		ON CONFLICT (user_id) DO UPDATE
		SET secret = EXCLUDED.secret,
		    created_at = NOW(),
		    expires_at = EXCLUDED.expires_at
	`
	_, err := r.db.ExecContext(ctx, query, pending.ID, pending.UserID, pending.Secret, pending.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create two factor pending: %w", err)
	}
	return nil
}

// GetTwoFactorPendingByUserID returns the pending record for a user
func (r *userRepository) GetTwoFactorPendingByUserID(ctx context.Context, userID uuid.UUID) (*TwoFactorPending, error) {
	query := `
		SELECT id, user_id, secret, created_at, expires_at
		FROM two_factor_pending
		WHERE user_id = $1
	`
	var pending TwoFactorPending
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&pending.ID, &pending.UserID, &pending.Secret, &pending.CreatedAt, &pending.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query two factor pending: %w", err)
	}
	return &pending, nil
}

// DeleteTwoFactorPending removes the pending record for a user
func (r *userRepository) DeleteTwoFactorPending(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM two_factor_pending WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete two factor pending: %w", err)
	}
	return nil
}

// GetByUsernameWithStats returns a user by username with their statistics
func (r *userRepository) GetByUsernameWithStats(ctx context.Context, username string) (*User, *UserStats, error) {
	query := `
		SELECT u.user_id, u.email, u.username, u.password_hash, u.phone_country_code, u.phone_number,
		       u.phone_verified, u.email_verified, u.status,
		       u.avatar_url, u.bio, u.timezone, u.locale, u.is_merchant, u.is_suspended,
		       u.suspension_reason, u.suspended_at, u.suspended_until, u.last_login_at, u.last_seen_at,
		       u.login_count, u.failed_login_attempts, u.locked_until, u.referral_code,
		       u.referred_by, u.two_fa_secret, u.remember_2fa, u.remember_2fa_expiry,
		       u.created_at, u.updated_at, u.deleted_at,
		       COALESCE(us.total_trades, 0) AS total_trades,
		       COALESCE(us.successful_trades, 0) AS successful_trades,
		       COALESCE(us.dispute_rate, 0) AS dispute_rate,
		       COALESCE(us.avg_rating, 0) AS avg_rating,
		       COALESCE(us.positive_feedback_count, 0) AS positive_feedback_count,
		       COALESCE(us.neutral_feedback_count, 0) AS neutral_feedback_count,
		       COALESCE(us.negative_feedback_count, 0) AS negative_feedback_count,
		       COALESCE(us.total_volume_usd, 0) AS total_volume_usd,
		       us.last_trade_at, us.updated_at AS stats_updated_at
		FROM users u
		LEFT JOIN user_stats us ON u.user_id = us.user_id
		WHERE u.username = $1 AND u.deleted_at IS NULL
	`
	var u User
	var stats UserStats
	var referredBy sql.NullString
	// stats timestamps are NULL when user has no stats row (LEFT JOIN)
	var lastTradeAt sql.NullTime
	var statsUpdatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.UserID, &u.Email, &u.Username, &u.PasswordHash,
		&u.PhoneCountryCode, &u.PhoneNumber, &u.PhoneVerified, &u.EmailVerified,
		&u.Status, &u.AvatarURL, &u.Bio,
		&u.Timezone, &u.Locale, &u.IsMerchant, &u.IsSuspended,
		&u.SuspensionReason, &u.SuspendedAt, &u.SuspendedUntil, &u.LastLoginAt, &u.LastSeenAt,
		&u.LoginCount, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.ReferralCode, &referredBy, &u.TwoFASecret, &u.Remember2FA,
		&u.Remember2FAExpiry, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		&stats.TotalTrades, &stats.SuccessfulTrades, &stats.DisputeRate,
		&stats.AvgRating, &stats.PositiveFeedbackCount, &stats.NeutralFeedbackCount,
		&stats.NegativeFeedbackCount, &stats.TotalVolumeUSD,
		&lastTradeAt, &statsUpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("query user by username with stats: %w", err)
	}
	if referredBy.Valid {
		parsed, err := uuid.Parse(referredBy.String)
		if err == nil {
			u.ReferredBy = NullUUID{UUID: parsed, Valid: true}
		}
	}
	if lastTradeAt.Valid {
		stats.LastTradeAt = &lastTradeAt.Time
	}
	if statsUpdatedAt.Valid {
		stats.UpdatedAt = statsUpdatedAt.Time
	}
	return &u, &stats, nil
}

// UpdateLastSeen updates the user's last_seen_at timestamp
func (r *userRepository) UpdateLastSeen(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET last_seen_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("update last seen: %w", err)
	}
	return nil
}

// BlockUser blocks a user
func (r *userRepository) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID, reason string, isPermanent bool, expiresAt *time.Time) error {
	query := `
		INSERT INTO user_blocks (blocker_id, blocked_id, reason, is_permanent, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (blocker_id, blocked_id) DO UPDATE SET
			reason = EXCLUDED.reason,
			is_permanent = EXCLUDED.is_permanent,
			expires_at = EXCLUDED.expires_at,
			created_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, blockerID, blockedID, reason, isPermanent, expiresAt)
	if err != nil {
		return fmt.Errorf("block user: %w", err)
	}
	return nil
}

// UnblockUser removes a user block
func (r *userRepository) UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	query := `DELETE FROM user_blocks WHERE blocker_id = $1 AND blocked_id = $2`
	_, err := r.db.ExecContext(ctx, query, blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("unblock user: %w", err)
	}
	return nil
}

// IsBlocked checks if a user is blocked by another user
func (r *userRepository) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_blocks
			WHERE blocker_id = $1 AND blocked_id = $2
			AND (is_permanent = true OR expires_at > NOW() OR expires_at IS NULL)
		)
	`
	var blocked bool
	err := r.db.QueryRowContext(ctx, query, blockerID, blockedID).Scan(&blocked)
	if err != nil {
		return false, fmt.Errorf("check blocked: %w", err)
	}
	return blocked, nil
}

// ListBlocks returns all users blocked by the given user
func (r *userRepository) ListBlocks(ctx context.Context, userID uuid.UUID) ([]UserBlock, error) {
	query := `
		SELECT id, blocker_id, blocked_id, reason, is_permanent, expires_at, created_at
		FROM user_blocks
		WHERE blocker_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query blocks: %w", err)
	}
	defer rows.Close()

	var blocks []UserBlock
	for rows.Next() {
		var b UserBlock
		err := rows.Scan(&b.ID, &b.BlockerID, &b.BlockedID, &b.Reason, &b.IsPermanent, &b.ExpiresAt, &b.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan block: %w", err)
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}
