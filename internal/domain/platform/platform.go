package platform

import (
	"time"

	"github.com/google/uuid"
)

// PaymentCategory represents payment method categories
type PaymentCategory string

const (
	PaymentCategoryMobileMoney  PaymentCategory = "mobile_money"
	PaymentCategoryBankTransfer PaymentCategory = "bank_transfer"
	PaymentCategoryOnlineWallet PaymentCategory = "online_wallet"
	PaymentCategoryCrypto       PaymentCategory = "crypto"
	PaymentCategoryCash         PaymentCategory = "cash"
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
	TargetAudience *string    `db:"target_audience" json:"target_audience,omitempty"` // all, verified_users
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

// IsForVerifiedUsers checks if the announcement targets verified users
func (a *Announcement) IsForVerifiedUsers() bool {
	return a.TargetAudience != nil && *a.TargetAudience == "verified_users"
}

// CryptoAsset represents a supported cryptocurrency
type CryptoAsset struct {
	ID              int       `db:"id" json:"id"`
	Symbol          string    `db:"symbol" json:"symbol"`
	Name            string    `db:"name" json:"name"`
	Blockchain      string    `db:"blockchain" json:"blockchain"`
	ContractAddress *string   `db:"contract_address" json:"contract_address,omitempty"`
	Decimals        int       `db:"decimals" json:"decimals"`
	IsActive        bool      `db:"is_active" json:"is_active"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// FiatCurrency represents a supported fiat currency
type FiatCurrency struct {
	ID        int       `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	Name      string    `db:"name" json:"name"`
	Symbol    string    `db:"symbol" json:"symbol"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// PaymentMethod represents a supported payment method
type PaymentMethod struct {
	ID          int             `db:"id" json:"id"`
	Code        string          `db:"code" json:"code"`
	Name        string          `db:"name" json:"name"`
	Category    PaymentCategory `db:"category" json:"category"`
	IconURL     *string         `db:"icon_url" json:"icon_url,omitempty"`
	Description *string         `db:"description" json:"description,omitempty"`
	IsActive    bool            `db:"is_active" json:"is_active"`
	SortOrder   int             `db:"sort_order" json:"sort_order"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
}
