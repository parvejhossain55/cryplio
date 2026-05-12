package dto

// UserStatsDTO represents trade statistics for a user
type UserStatsDTO struct {
	TotalTrades           int      `json:"total_trades"`
	SuccessfulTrades      int      `json:"successful_trades"`
	DisputeRate           float64  `json:"dispute_rate"`
	AvgRating             *float64 `json:"avg_rating,omitempty"`
	PositiveFeedbackCount int      `json:"positive_feedback_count"`
	NeutralFeedbackCount  int      `json:"neutral_feedback_count"`
	NegativeFeedbackCount int      `json:"negative_feedback_count"`
	TotalVolumeUSD        float64  `json:"total_volume_usd"`
	LastTradeAt           string   `json:"last_trade_at,omitempty"`
}

// UserResponse is the delivery-safe projection of a user.
type UserResponse struct {
	ID            string       `json:"id"`
	Email         string       `json:"email"`
	Username      string       `json:"username"`
	EmailVerified bool         `json:"email_verified"`
	IsMerchant    bool         `json:"is_merchant"`
	TwoFAEnabled  bool         `json:"two_fa_enabled"`
	AvatarURL     *string      `json:"avatar_url,omitempty"`
	Bio           *string      `json:"bio,omitempty"`
	LastSeenAt    string       `json:"last_seen_at,omitempty"`
	IsOnline      bool         `json:"is_online"`
	Stats         UserStatsDTO `json:"stats"`
	// Header profile fields
	TraderBadge             string `json:"trader_badge,omitempty"` // e.g., "PRO TRADER", "VERIFIED", ""
	UnreadNotificationCount int    `json:"unread_notification_count,omitempty"`
	AccountHealth           string `json:"account_health,omitempty"`      // e.g., "EXCELLENT", "GOOD", "FAIR", "POOR"
	AccountSecurity         string `json:"account_security,omitempty"`    // e.g., "VERIFIED", "UNVERIFIED"
	TwoFactorStatus         string `json:"two_factor_status,omitempty"`   // e.g., "ENABLED", "DISABLED"
	LoginNotifications      string `json:"login_notifications,omitempty"` // e.g., "ACTIVE", "INACTIVE"
}

// AuthResponse is returned from login and registration flows.
type AuthResponse struct {
	Token        string       `json:"token"`                   // access token
	RefreshToken string       `json:"refresh_token,omitempty"` // refresh token
	User         UserResponse `json:"user"`
}
