package dto

// UserResponse is the delivery-safe projection of a user.
type UserResponse struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	Username      string  `json:"username"`
	EmailVerified bool    `json:"email_verified"`
	KYCLevel      int     `json:"kyc_level"`
	IsMerchant    bool    `json:"is_merchant"`
	TwoFAEnabled  bool    `json:"two_fa_enabled"`
	AvatarURL     *string `json:"avatar_url,omitempty"`
	Bio           *string `json:"bio,omitempty"`
}

// AuthResponse is returned from login and registration flows.
type AuthResponse struct {
	Token        string       `json:"token"`                   // access token
	RefreshToken string       `json:"refresh_token,omitempty"` // refresh token
	User         UserResponse `json:"user"`
}
