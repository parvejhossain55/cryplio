package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// ─── Core User CRUD ───────────────────────────────────────────────────────────

// GetAll returns all non-deleted users ordered by creation date (admin use).
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := scanUser(rows, &u); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

// GetByID returns a single non-deleted user by primary key.
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE user_id = $1 AND deleted_at IS NULL`
	var u User
	if err := scanUser(r.db.QueryRowContext(ctx, query, id), &u); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	return &u, nil
}

// GetByEmail returns a single non-deleted user by email address.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE email = $1 AND deleted_at IS NULL`
	var u User
	if err := scanUser(r.db.QueryRowContext(ctx, query, email), &u); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by email: %w", err)
	}
	return &u, nil
}

// GetByUsername returns a single non-deleted user by username.
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE username = $1 AND deleted_at IS NULL`
	var u User
	if err := scanUser(r.db.QueryRowContext(ctx, query, username), &u); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by username: %w", err)
	}
	return &u, nil
}

// GetByUsernameWithStats returns a user together with their aggregated trade
// statistics in a single LEFT JOIN query.
func (r *userRepository) GetByUsernameWithStats(ctx context.Context, username string) (*User, *UserStats, error) {
	query := `
		SELECT
			u.user_id, u.email, u.username, u.password_hash,
			u.phone_country_code, u.phone_number, u.phone_verified, u.email_verified,
			u.status, u.role,
			u.avatar_url, u.bio, u.timezone, u.locale,
			u.is_merchant, u.is_suspended,
			u.suspension_reason, u.suspended_at, u.suspended_until,
			u.last_login_at, u.last_seen_at,
			u.login_count, u.failed_login_attempts, u.locked_until,
			u.referral_code, u.referred_by,
			u.two_fa_secret, u.remember_2fa, u.remember_2fa_expiry,
			u.created_at, u.updated_at, u.deleted_at,
			COALESCE(s.total_trades, 0),
			COALESCE(s.successful_trades, 0),
			COALESCE(s.dispute_rate, 0),
			COALESCE(s.avg_rating, 0),
			COALESCE(s.positive_feedback_count, 0),
			COALESCE(s.neutral_feedback_count, 0),
			COALESCE(s.negative_feedback_count, 0),
			COALESCE(s.total_volume_usd, 0),
			s.last_trade_at,
			s.updated_at
		FROM users u
		LEFT JOIN user_stats s ON u.user_id = s.user_id
		WHERE u.username = $1 AND u.deleted_at IS NULL
	`
	var u User
	var stats UserStats
	var lastTradeAt sql.NullTime
	var statsUpdatedAt sql.NullTime

	row := r.db.QueryRowContext(ctx, query, username)
	var referredBy sql.NullString
	err := row.Scan(
		&u.UserID, &u.Email, &u.Username, &u.PasswordHash,
		&u.PhoneCountryCode, &u.PhoneNumber, &u.PhoneVerified, &u.EmailVerified,
		&u.Status, &u.Role,
		&u.AvatarURL, &u.Bio, &u.Timezone, &u.Locale,
		&u.IsMerchant, &u.IsSuspended,
		&u.SuspensionReason, &u.SuspendedAt, &u.SuspendedUntil,
		&u.LastLoginAt, &u.LastSeenAt,
		&u.LoginCount, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.ReferralCode, &referredBy,
		&u.TwoFASecret, &u.Remember2FA, &u.Remember2FAExpiry,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		&stats.TotalTrades, &stats.SuccessfulTrades, &stats.DisputeRate,
		&stats.AvgRating,
		&stats.PositiveFeedbackCount, &stats.NeutralFeedbackCount, &stats.NegativeFeedbackCount,
		&stats.TotalVolumeUSD,
		&lastTradeAt, &statsUpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("query user by username with stats: %w", err)
	}
	if referredBy.Valid {
		if parsed, err := uuid.Parse(referredBy.String); err == nil {
			u.ReferredBy = NullUUID{UUID: parsed, Valid: true}
		}
	}
	if lastTradeAt.Valid {
		stats.LastTradeAt = &lastTradeAt.Time
	}
	if statsUpdatedAt.Valid {
		stats.UpdatedAt = statsUpdatedAt.Time
	}
	stats.UserID = u.UserID
	return &u, &stats, nil
}

// ─── Mutations ────────────────────────────────────────────────────────────────

// Create inserts a new user row and updates CreatedAt / UpdatedAt from the DB.
func (r *userRepository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (
			user_id, email, username, password_hash,
			phone_country_code, phone_number, phone_verified, email_verified,
			status, role, timezone, locale,
			is_merchant, login_count, failed_login_attempts,
			referral_code, remember_2fa, referred_by,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18,
			NOW(), NOW()
		) RETURNING created_at, updated_at
	`
	var referredBy sql.NullString
	if u.ReferredBy.Valid {
		referredBy = sql.NullString{String: u.ReferredBy.UUID.String(), Valid: true}
	}
	return r.db.QueryRowContext(ctx, query,
		u.UserID, u.Email, u.Username, u.PasswordHash,
		u.PhoneCountryCode, u.PhoneNumber, u.PhoneVerified, u.EmailVerified,
		u.Status, u.Role, u.Timezone, u.Locale,
		u.IsMerchant, u.LoginCount, u.FailedLoginAttempts,
		u.ReferralCode, u.Remember2FA, referredBy,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
}

// Update persists all mutable user fields and refreshes UpdatedAt.
func (r *userRepository) Update(ctx context.Context, u *User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2,
		    phone_country_code = $3, phone_number = $4,
		    phone_verified = $5, email_verified = $6,
		    avatar_url = $7, bio = $8, timezone = $9, locale = $10,
		    is_merchant = $11, is_suspended = $12,
		    suspension_reason = $13, suspended_at = $14, suspended_until = $15,
		    last_login_at = $16, login_count = $17, failed_login_attempts = $18,
		    locked_until = $19, referral_code = $20, referred_by = $21,
		    two_fa_secret = $22, remember_2fa = $23, remember_2fa_expiry = $24,
		    updated_at = NOW()
		WHERE user_id = $25 AND deleted_at IS NULL
		RETURNING updated_at
	`
	var referredBy sql.NullString
	if u.ReferredBy.Valid {
		referredBy = sql.NullString{String: u.ReferredBy.UUID.String(), Valid: true}
	}
	err := r.db.QueryRowContext(ctx, query,
		u.Email, u.Username,
		u.PhoneCountryCode, u.PhoneNumber,
		u.PhoneVerified, u.EmailVerified,
		u.AvatarURL, u.Bio, u.Timezone, u.Locale,
		u.IsMerchant, u.IsSuspended,
		u.SuspensionReason, u.SuspendedAt, u.SuspendedUntil,
		u.LastLoginAt, u.LoginCount, u.FailedLoginAttempts,
		u.LockedUntil, u.ReferralCode, referredBy,
		u.TwoFASecret, u.Remember2FA, u.Remember2FAExpiry,
		u.UserID,
	).Scan(&u.UpdatedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found")
	}
	return err
}

// Delete soft-deletes the user by setting deleted_at.
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE users SET deleted_at = NOW() WHERE user_id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}

// ─── Login tracking ───────────────────────────────────────────────────────────

// IncrementLogin records a successful login: bumps login_count, resets the
// failed-attempt counter and any lock, and sets last_login_at to NOW().
func (r *userRepository) IncrementLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET last_login_at = NOW(),
		    login_count = login_count + 1,
		    failed_login_attempts = 0,
		    locked_until = NULL,
		    updated_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL`, userID)
	if err != nil {
		return fmt.Errorf("increment login: %w", err)
	}
	return nil
}

// IncrementFailedAttempts bumps the counter and auto-locks the account after
// 5 consecutive failures (lock duration: 15 minutes).
func (r *userRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error) {
	var attempts int
	err := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET failed_login_attempts = failed_login_attempts + 1,
		    locked_until = CASE
		        WHEN failed_login_attempts >= 4 THEN NOW() + INTERVAL '15 minutes'
		        ELSE locked_until
		    END,
		    updated_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL
		RETURNING failed_login_attempts`, userID).Scan(&attempts)
	if err != nil {
		return 0, fmt.Errorf("increment failed attempts: %w", err)
	}
	return attempts, nil
}

// UpdateLastSeen refreshes last_seen_at to the current timestamp.
func (r *userRepository) UpdateLastSeen(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET last_seen_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL`, userID)
	return err
}

// ─── Aggregates ───────────────────────────────────────────────────────────────

// CountUsers returns the total number of non-deleted users.
func (r *userRepository) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}

// GetStats returns the pre-aggregated trade statistics for a user.
// Returns nil (no error) when the user has no stats row yet.
func (r *userRepository) GetStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	query := `
		SELECT user_id, total_trades, successful_trades, dispute_rate, avg_rating,
		       positive_feedback_count, neutral_feedback_count, negative_feedback_count,
		       total_volume_usd, last_trade_at, updated_at
		FROM user_stats
		WHERE user_id = $1
	`
	var s UserStats
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&s.UserID, &s.TotalTrades, &s.SuccessfulTrades, &s.DisputeRate, &s.AvgRating,
		&s.PositiveFeedbackCount, &s.NeutralFeedbackCount, &s.NegativeFeedbackCount,
		&s.TotalVolumeUSD, &s.LastTradeAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query user stats: %w", err)
	}
	return &s, nil
}
