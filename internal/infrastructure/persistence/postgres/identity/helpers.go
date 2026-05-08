package identity

import (
	"database/sql"

	"github.com/google/uuid"
)

// ─── Shared scan helpers ──────────────────────────────────────────────────────

// scanner is satisfied by both *sql.Row and *sql.Rows, letting scanUser work
// for both single-row queries and cursor-based loops.
type scanner interface {
	Scan(dest ...any) error
}

// userColumns is the canonical SELECT column list for the users table.
// Every query that returns a full User must use this list so that scanUser
// can scan the result in a single call.
//
// FIXED: adds `role` and `last_seen_at` which were missing in the original
// single-file implementation (role was always empty after a DB read,
// breaking admin JWT role claims).
const userColumns = `
	user_id, email, username, password_hash,
	phone_country_code, phone_number, phone_verified, email_verified,
	status, role,
	avatar_url, bio, timezone, locale,
	is_merchant, is_suspended,
	suspension_reason, suspended_at, suspended_until,
	last_login_at, last_seen_at,
	login_count, failed_login_attempts, locked_until,
	referral_code, referred_by,
	two_fa_secret, remember_2fa, remember_2fa_expiry,
	created_at, updated_at, deleted_at
`

// scanUser reads the 32-column user projection (userColumns) into u.
// It handles the referred_by UUID nullable conversion internally.
func scanUser(row scanner, u *User) error {
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
	)
	if err != nil {
		return err
	}
	if referredBy.Valid {
		if parsed, err := uuid.Parse(referredBy.String); err == nil {
			u.ReferredBy = NullUUID{UUID: parsed, Valid: true}
		}
	}
	return nil
}
