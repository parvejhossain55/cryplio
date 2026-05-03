package events

import "time"

// UserRegisteredEvent is raised when a new user registers.
type UserRegisteredEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (e UserRegisteredEvent) Name() string { return "user.registered" }

// UserLoggedInEvent is raised when a user logs in.
type UserLoggedInEvent struct {
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	LoggedInAt time.Time `json:"logged_in_at"`
}

func (e UserLoggedInEvent) Name() string { return "user.logged_in" }

// UserLoggedOutEvent is raised when a user logs out.
type UserLoggedOutEvent struct {
	UserID      string    `json:"user_id"`
	LoggedOutAt time.Time `json:"logged_out_at"`
}

func (e UserLoggedOutEvent) Name() string { return "user.logged_out" }

// EmailVerificationRequestedEvent is raised when email verification is requested.
type EmailVerificationRequestedEvent struct {
	UserID  string    `json:"user_id"`
	Email   string    `json:"email"`
	Token   string    `json:"-"` // token sent to user, not for logging
	Expires time.Time `json:"expires_at"`
}

func (e EmailVerificationRequestedEvent) Name() string { return "user.email_verification_requested" }

// EmailVerifiedEvent is raised when email is verified.
type EmailVerifiedEvent struct {
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	VerifiedAt time.Time `json:"verified_at"`
}

func (e EmailVerifiedEvent) Name() string { return "user.email_verified" }

// PasswordResetRequestedEvent is raised when a password reset is requested.
type PasswordResetRequestedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Token  string `json:"-"` // token sent to user
}

func (e PasswordResetRequestedEvent) Name() string { return "user.password_reset_requested" }

// PasswordResetCompletedEvent is raised when a password is reset.
type PasswordResetCompletedEvent struct {
	UserID  string    `json:"user_id"`
	Email   string    `json:"email"`
	ResetAt time.Time `json:"reset_at"`
}

func (e PasswordResetCompletedEvent) Name() string { return "user.password_reset_completed" }

// TwoFactorEnabledEvent is raised when 2FA is enabled.
type TwoFactorEnabledEvent struct {
	UserID    string    `json:"user_id"`
	EnabledAt time.Time `json:"enabled_at"`
}

func (e TwoFactorEnabledEvent) Name() string { return "user.two_factor_enabled" }
