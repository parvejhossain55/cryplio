package identity

import (
	"time"

	"github.com/google/uuid"
)

// UserOAuth represents an OAuth connection for a user
type UserOAuth struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	UserID           uuid.UUID  `db:"user_id" json:"user_id"`
	Provider         string     `db:"provider" json:"provider"`
	ProviderUserID   string     `db:"provider_user_id" json:"provider_user_id"`
	ProviderEmail    *string    `db:"provider_email" json:"provider_email,omitempty"`
	ProviderUsername *string    `db:"provider_username" json:"provider_username,omitempty"`
	AccessToken      *string    `db:"access_token" json:"access_token,omitempty"`
	RefreshToken     *string    `db:"refresh_token" json:"refresh_token,omitempty"`
	TokenExpiry      *time.Time `db:"token_expiry" json:"token_expiry,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// IsTokenExpired checks if the OAuth token is expired
func (o *UserOAuth) IsTokenExpired() bool {
	if o.TokenExpiry == nil {
		return true
	}
	return time.Now().After(*o.TokenExpiry)
}

// NeedsRefresh checks if the token needs to be refreshed
func (o *UserOAuth) NeedsRefresh() bool {
	if o.TokenExpiry == nil || o.RefreshToken == nil {
		return false
	}
	// Refresh 5 minutes before expiry
	refreshThreshold := *o.TokenExpiry
	refreshThreshold = refreshThreshold.Add(-5 * time.Minute)
	return time.Now().After(refreshThreshold)
}
