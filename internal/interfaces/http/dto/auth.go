package dto

import "github.com/google/uuid"

// RegisterRequest is the auth registration payload.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest is the auth login payload.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateProfileRequest is the authenticated profile update payload.
type UpdateProfileRequest struct {
	Username *string `json:"username,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

// EmailVerificationRequest requests a verification email for a user (admin or self).
type EmailVerificationRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

// EmailVerifyRequest verifies email using token.
type EmailVerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

// PasswordResetRequest initiates a password reset.
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirm confirms password reset with token and new password.
type PasswordResetConfirm struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// RefreshTokenRequest requests new access token using a refresh token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TwoFactorSetupResponse returns secret and provisioning URI for TOTP.
type TwoFactorSetupResponse struct {
	Secret          string `json:"secret"`
	ProvisioningURI string `json:"provisioning_uri"`
}

// TwoFactorVerifyRequest verifies and enables 2FA.
type TwoFactorVerifyRequest struct {
	Code string `json:"code" binding:"required"`
}

// TwoFactorDisableRequest disables 2FA after password confirmation.
type TwoFactorDisableRequest struct {
	Password string `json:"password" binding:"required"`
}

// TwoFactorLoginCompleteRequest completes the 2FA login step.
type TwoFactorLoginCompleteRequest struct {
	TempToken string `json:"temp_token" binding:"required"`
	Code      string `json:"code" binding:"required"`
}
