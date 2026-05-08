package identity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusPending   UserStatus = "pending"
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
	UserStatusDeleted   UserStatus = "deleted"
)

// UserRole represents the user role
type UserRole string

const (
	UserRoleUser     UserRole = "user"
	UserRoleMerchant UserRole = "merchant"
	UserRoleAdmin    UserRole = "admin"
)

// NullUUID wraps uuid.UUID to handle NULL values
type NullUUID struct {
	UUID  uuid.UUID
	Valid bool
}

// Scan implements sql.Scanner interface
func (n *NullUUID) Scan(value interface{}) error {
	if value == nil {
		n.UUID = uuid.Nil
		n.Valid = false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		u, err := uuid.FromBytes(v)
		if err != nil {
			return err
		}
		n.UUID = u
		n.Valid = true
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		n.UUID = u
		n.Valid = true
	case uuid.UUID:
		n.UUID = v
		n.Valid = true
	}
	return nil
}

// Value implements driver.Valuer interface
func (n NullUUID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.UUID[:], nil
}

// MarshalJSON implements json.Marshaler
func (n NullUUID) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.UUID.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (n *NullUUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.UUID = uuid.Nil
		n.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	n.UUID = u
	n.Valid = true
	return nil
}

// User represents a user account
type User struct {
	UserID              uuid.UUID  `db:"user_id" json:"user_id"`
	Email               string     `db:"email" json:"email"`
	Username            string     `db:"username" json:"username"`
	PasswordHash        string     `db:"password_hash" json:"-"`
	PhoneCountryCode    *string    `db:"phone_country_code" json:"phone_country_code,omitempty"`
	PhoneNumber         *string    `db:"phone_number" json:"phone_number,omitempty"`
	PhoneVerified       bool       `db:"phone_verified" json:"phone_verified"`
	EmailVerified       bool       `db:"email_verified" json:"email_verified"`
	Status              UserStatus `db:"status" json:"status"`
	Role                UserRole   `db:"role" json:"role"`
	AvatarURL           *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	Bio                 *string    `db:"bio" json:"bio,omitempty"`
	Timezone            string     `db:"timezone" json:"timezone"`
	Locale              string     `db:"locale" json:"locale"`
	IsMerchant          bool       `db:"is_merchant" json:"is_merchant"`
	IsSuspended         bool       `db:"is_suspended" json:"is_suspended"`
	SuspensionReason    *string    `db:"suspension_reason" json:"suspension_reason,omitempty"`
	SuspendedAt         *time.Time `db:"suspended_at" json:"suspended_at,omitempty"`
	SuspendedUntil      *time.Time `db:"suspended_until" json:"suspended_until,omitempty"`
	UnsuspendedAt       *time.Time `db:"unsuspended_at" json:"unsuspended_at,omitempty"`
	LastLoginAt         *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	LastSeenAt          *time.Time `db:"last_seen_at" json:"last_seen_at,omitempty"`
	LoginCount          int        `db:"login_count" json:"login_count"`
	FailedLoginAttempts int        `db:"failed_login_attempts" json:"failed_login_attempts"`
	LockedUntil         *time.Time `db:"locked_until" json:"locked_until,omitempty"`
	ReferralCode        *string    `db:"referral_code" json:"referral_code,omitempty"`
	ReferredBy          NullUUID   `db:"referred_by" json:"referred_by,omitempty"`
	TwoFASecret         *string    `db:"two_fa_secret" json:"-"`
	Remember2FA         bool       `db:"remember_2fa" json:"remember_2fa"`
	Remember2FAExpiry   *time.Time `db:"remember_2fa_expiry" json:"remember_2fa_expiry,omitempty"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// IsOnline checks if the user was seen in the last 5 minutes
func (u *User) IsOnline() bool {
	if u.LastSeenAt == nil {
		return false
	}
	return time.Since(*u.LastSeenAt) < 5*time.Minute
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && u.DeletedAt == nil && !u.IsSuspended
}

// IsDeleted checks if the user is soft-deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// SetDeleted marks the user as soft-deleted
func (u *User) SetDeleted(deletedAt time.Time) {
	u.DeletedAt = &deletedAt
	u.Status = UserStatusDeleted
}

// IncrementLoginCount increments the login count
func (u *User) IncrementLoginCount() {
	u.LoginCount++
}

// IncrementFailedAttempts increments failed login attempts
func (u *User) IncrementFailedAttempts() {
	u.FailedLoginAttempts++
}

// LockAccount locks the user account until the given time
func (u *User) LockAccount(until time.Time) {
	u.LockedUntil = &until
	u.FailedLoginAttempts = 0
}

// UnlockAccount unlocks the user account
func (u *User) UnlockAccount() {
	u.LockedUntil = nil
	u.FailedLoginAttempts = 0
}

// IsLocked checks if the account is currently locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// Suspend suspends the user account
func (u *User) Suspend(reason string, until *time.Time) {
	now := time.Now()
	u.IsSuspended = true
	u.Status = UserStatusSuspended
	u.SuspensionReason = &reason
	u.SuspendedAt = &now
	u.SuspendedUntil = until
	u.UpdatedAt = now
}

// Unsuspend lifts the suspension on the user account
func (u *User) Unsuspend() {
	now := time.Now()
	u.IsSuspended = false
	u.Status = UserStatusActive
	u.SuspensionReason = nil
	u.SuspendedUntil = nil
	u.UnsuspendedAt = &now
	u.UpdatedAt = now
}

// NewUser creates a new user
func NewUser(email, username, passwordHash string) *User {
	now := time.Now()
	return &User{
		UserID:              uuid.New(),
		Email:               email,
		Username:            username,
		PasswordHash:        passwordHash,
		Status:              UserStatusActive,
		Role:                UserRoleUser,
		Timezone:            "UTC",
		Locale:              "en",
		PhoneVerified:       false,
		EmailVerified:       false,
		IsMerchant:          false,
		IsSuspended:         false,
		LoginCount:          0,
		FailedLoginAttempts: 0,
		Remember2FA:         false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// UserSession represents a user login session
type UserSession struct {
	ID                uuid.UUID `db:"id" json:"id"`
	UserID            uuid.UUID `db:"user_id" json:"user_id"`
	TokenID           string    `db:"token_id" json:"token_id"`
	DeviceFingerprint *string   `db:"device_fingerprint" json:"device_fingerprint,omitempty"`
	IPAddress         *string   `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent         *string   `db:"user_agent" json:"user_agent,omitempty"`
	DeviceType        *string   `db:"device_type" json:"device_type,omitempty"`
	Location          *string   `db:"location" json:"location,omitempty"`
	IsRemembered      bool      `db:"is_remembered" json:"is_remembered"`
	ExpiresAt         time.Time `db:"expires_at" json:"expires_at"`
	LastUsedAt        time.Time `db:"last_used_at" json:"last_used_at"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}

// IsExpired checks if the session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// PasswordResetToken stores password reset tokens
type PasswordResetToken struct {
	ID        int        `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	TokenHash string     `db:"token_hash" json:"-"`
	IPAddress *string    `db:"ip_address" json:"ip_address,omitempty"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	UsedAt    *time.Time `db:"used_at" json:"used_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

// IsExpired checks if the token has expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has been used
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// EmailVerificationToken stores email verification tokens
type EmailVerificationToken struct {
	ID         int        `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	TokenHash  string     `db:"token_hash" json:"-"`
	ExpiresAt  time.Time  `db:"expires_at" json:"expires_at"`
	VerifiedAt *time.Time `db:"verified_at" json:"verified_at,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

// IsExpired checks if the token has expired
func (t *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsVerified checks if the email has been verified
func (t *EmailVerificationToken) IsVerified() bool {
	return t.VerifiedAt != nil
}

// UserBlock represents a user blocking relationship
type UserBlock struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	BlockerID   uuid.UUID  `db:"blocker_id" json:"blocker_id"`
	BlockedID   uuid.UUID  `db:"blocked_id" json:"blocked_id"`
	Reason      *string    `db:"reason" json:"reason,omitempty"`
	IsPermanent bool       `db:"is_permanent" json:"is_permanent"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

// IsActive checks if the block is currently active
func (b *UserBlock) IsActive() bool {
	if b.IsPermanent {
		return true
	}
	if b.ExpiresAt == nil {
		return true
	}
	return time.Now().Before(*b.ExpiresAt)
}

// UserStats represents denormalized user statistics
type UserStats struct {
	UserID                uuid.UUID  `db:"user_id" json:"user_id"`
	TotalTrades           int        `db:"total_trades" json:"total_trades"`
	SuccessfulTrades      int        `db:"successful_trades" json:"successful_trades"`
	DisputeRate           float64    `db:"dispute_rate" json:"dispute_rate"`
	AvgRating             *float64   `db:"avg_rating" json:"avg_rating,omitempty"`
	PositiveFeedbackCount int        `db:"positive_feedback_count" json:"positive_feedback_count"`
	NeutralFeedbackCount  int        `db:"neutral_feedback_count" json:"neutral_feedback_count"`
	NegativeFeedbackCount int        `db:"negative_feedback_count" json:"negative_feedback_count"`
	TotalVolumeUSD        float64    `db:"total_volume_usd" json:"total_volume_usd"`
	LastTradeAt           *time.Time `db:"last_trade_at" json:"last_trade_at,omitempty"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updated_at"`
}

// GetSuccessRate returns the success rate as a percentage
func (s *UserStats) GetSuccessRate() float64 {
	if s.TotalTrades == 0 {
		return 0.0
	}
	return float64(s.SuccessfulTrades) / float64(s.TotalTrades) * 100.0
}
