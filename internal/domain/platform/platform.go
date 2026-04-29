package platform

import (
	"time"

	"github.com/google/uuid"
)

// PlatformConfig represents a key-value platform configuration
type PlatformConfig struct {
	Key         string    `db:"key" json:"key"`
	Value       string    `db:"value" json:"value"`
	Description *string   `db:"description" json:"description,omitempty"`
	UpdatedBy   *string   `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// IsEmpty checks if the config value is empty
func (p *PlatformConfig) IsEmpty() bool {
	return p.Value == ""
}

// GetString returns the config value as a string
func (p *PlatformConfig) GetString() string {
	return p.Value
}

// Announcement represents a system announcement
type Announcement struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	Title          string     `db:"title" json:"title"`
	Message        string     `db:"message" json:"message"`
	Type           string     `db:"type" json:"type"` // info, warning, critical
	IsActive       bool       `db:"is_active" json:"is_active"`
	TargetAudience *string    `db:"target_audience" json:"target_audience,omitempty"` // all, merchants, verified_users
	ExpiresAt      *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"created_by,omitempty"` // admin user ID
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

// IsExpired checks if the announcement has expired
func (a *Announcement) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// IsActiveNow checks if the announcement is currently active
func (a *Announcement) IsActiveNow() bool {
	return a.IsActive && (a.ExpiresAt == nil || time.Now().Before(*a.ExpiresAt))
}

// IsForAll checks if the announcement targets all users
func (a *Announcement) IsForAll() bool {
	return a.TargetAudience == nil || *a.TargetAudience == "all"
}

// IsForMerchants checks if the announcement targets merchants
func (a *Announcement) IsForMerchants() bool {
	return a.TargetAudience != nil && *a.TargetAudience == "merchants"
}

// IsForVerifiedUsers checks if the announcement targets verified users
func (a *Announcement) IsForVerifiedUsers() bool {
	return a.TargetAudience != nil && *a.TargetAudience == "verified_users"
}
