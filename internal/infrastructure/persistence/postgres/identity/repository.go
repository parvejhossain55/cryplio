package identity

import (
	"context"
	"database/sql"
	"fmt"

	domainidentity "cryplio/internal/domain/identity"
	"github.com/google/uuid"
)

type User = domainidentity.User
type UserStats = domainidentity.UserStats
type UserOAuth = domainidentity.UserOAuth
type NullUUID = domainidentity.NullUUID

// UserRepository defines data access for users
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, limit, offset int) ([]User, error)
	IncrementLogin(ctx context.Context, userID uuid.UUID) error
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error)
	GetStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
	// OAuth
	GetOAuthByProviderID(ctx context.Context, provider, providerUserID string) (*UserOAuth, error)
	CreateOAuth(ctx context.Context, oauth *UserOAuth) error
	UpdateOAuth(ctx context.Context, oauth *UserOAuth) error
	GetOAuthByUserID(ctx context.Context, userID uuid.UUID) ([]UserOAuth, error)
	DeleteOAuth(ctx context.Context, id uuid.UUID) error
}

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
		       phone_verified, email_verified, kyc_level, kyc_last_updated, status,
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
			&u.KYCLevel, &u.KYCLastUpdated, &u.Status, &u.AvatarURL, &u.Bio,
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
		       phone_verified, email_verified, kyc_level, kyc_last_updated, status,
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
		&u.KYCLevel, &u.KYCLastUpdated, &u.Status, &u.AvatarURL, &u.Bio,
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
		       phone_verified, email_verified, kyc_level, kyc_last_updated, status,
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
		&u.KYCLevel, &u.KYCLastUpdated, &u.Status, &u.AvatarURL, &u.Bio,
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
		       phone_verified, email_verified, kyc_level, kyc_last_updated, status,
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
		&u.KYCLevel, &u.KYCLastUpdated, &u.Status, &u.AvatarURL, &u.Bio,
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
			phone_verified, email_verified, kyc_level, status, timezone, locale,
			is_merchant, login_count, failed_login_attempts, referral_code, remember_2fa,
			referred_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
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
		u.KYCLevel, u.Status, u.Timezone, u.Locale, u.IsMerchant,
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
		    phone_verified = $5, email_verified = $6, kyc_level = $7,
		    avatar_url = $8, bio = $9, timezone = $10, locale = $11,
		    is_merchant = $12, is_suspended = $13, suspension_reason = $14,
		    suspended_at = $15, suspended_until = $16, last_login_at = $17,
		    login_count = $18, failed_login_attempts = $19, locked_until = $20,
		    referral_code = $21, referred_by = $22, two_fa_secret = $23,
		    remember_2fa = $24, remember_2fa_expiry = $25, updated_at = NOW()
		WHERE user_id = $26 AND deleted_at IS NULL
		RETURNING updated_at
	`
	nullReferredBy := sql.NullString{}
	if u.ReferredBy.Valid {
		nullReferredBy = sql.NullString{String: u.ReferredBy.UUID.String(), Valid: true}
	}
	err := r.db.QueryRowContext(
		ctx, query,
		u.Email, u.Username, u.PhoneCountryCode, u.PhoneNumber,
		u.PhoneVerified, u.EmailVerified, u.KYCLevel, u.AvatarURL, u.Bio,
		u.Timezone, u.Locale, u.IsMerchant, u.IsSuspended, u.SuspensionReason,
		u.SuspendedAt, u.SuspendedUntil, u.LastLoginAt, u.LoginCount,
		u.FailedLoginAttempts, u.LockedUntil, u.ReferralCode,
		nullReferredBy, u.TwoFASecret, u.Remember2FA, u.Remember2FAExpiry,
		u.UserID,
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
