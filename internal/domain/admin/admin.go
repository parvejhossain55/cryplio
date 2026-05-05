package admin

import (
	"time"

	"github.com/google/uuid"
)

// MerchantStatus represents merchant application status
type MerchantStatus string

const (
	MerchantStatusNone      MerchantStatus = "none"
	MerchantStatusPending   MerchantStatus = "pending"
	MerchantStatusApproved  MerchantStatus = "approved"
	MerchantStatusRejected  MerchantStatus = "rejected"
	MerchantStatusSuspended MerchantStatus = "suspended"
)

// MerchantApplication represents a merchant verification application
type MerchantApplication struct {
	ApplicationID  uuid.UUID      `db:"application_id" json:"application_id"`
	UserID         uuid.UUID      `db:"user_id" json:"user_id"`
	BusinessName   string         `db:"business_name" json:"business_name"`
	BusinessType   string         `db:"business_type" json:"business_type"`
	TaxID          *string        `db:"tax_id" json:"tax_id,omitempty"`
	Website        *string        `db:"website" json:"website,omitempty"`
	Documents      []string       `db:"documents" json:"documents"` // JSONB array of document URLs
	Status         MerchantStatus `db:"status" json:"status"`
	RejectedReason *string        `db:"rejected_reason" json:"rejected_reason,omitempty"`
	ReviewedBy     *uuid.UUID     `db:"reviewed_by" json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time     `db:"reviewed_at" json:"reviewed_at,omitempty"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// IsPending checks if the application is pending review
func (m *MerchantApplication) IsPending() bool {
	return m.Status == MerchantStatusPending
}

// IsApproved checks if the application is approved
func (m *MerchantApplication) IsApproved() bool {
	return m.Status == MerchantStatusApproved
}

// IsRejected checks if the application is rejected
func (m *MerchantApplication) IsRejected() bool {
	return m.Status == MerchantStatusRejected
}

// Approve approves the merchant application
func (m *MerchantApplication) Approve(adminID uuid.UUID) {
	m.Status = MerchantStatusApproved
	now := time.Now()
	m.ReviewedBy = &adminID
	m.ReviewedAt = &now
}

// Reject rejects the merchant application with a reason
func (m *MerchantApplication) Reject(adminID uuid.UUID, reason string) {
	m.Status = MerchantStatusRejected
	now := time.Now()
	m.ReviewedBy = &adminID
	m.ReviewedAt = &now
	m.RejectedReason = &reason
}

// Suspend suspends the merchant
func (m *MerchantApplication) Suspend() {
	m.Status = MerchantStatusSuspended
}

// MerchantAnalytics represents daily merchant metrics
type MerchantAnalytics struct {
	ID            uuid.UUID `db:"id" json:"id"`
	MerchantID    uuid.UUID `db:"merchant_id" json:"merchant_id"`
	Date          time.Time `db:"date" json:"date"`
	TotalSales    float64   `db:"total_sales" json:"total_sales"`
	TotalVolume   float64   `db:"total_volume" json:"total_volume"`
	TradeCount    int       `db:"trade_count" json:"trade_count"`
	CustomerCount int       `db:"customer_count" json:"customer_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// IsEmpty checks if analytics has no data
func (m *MerchantAnalytics) IsEmpty() bool {
	return m.TradeCount == 0 && m.TotalSales == 0
}

// AdminActionType represents types of admin actions
type AdminActionType string

const (
	AdminActionUserSuspend       AdminActionType = "user_suspend"
	AdminActionUserBan           AdminActionType = "user_ban"
	AdminActionUserUnban         AdminActionType = "user_unban"
	AdminActionDisputeResolve    AdminActionType = "dispute_resolve"
	AdminActionWithdrawalApprove AdminActionType = "withdrawal_approve"
	AdminActionWithdrawalReject  AdminActionType = "withdrawal_reject"
	AdminActionAnnouncementPost  AdminActionType = "announcement_post"
	AdminActionFeeUpdate         AdminActionType = "fee_update"
	AdminActionMerchantApprove   AdminActionType = "merchant_approve"
	AdminActionMerchantReject    AdminActionType = "merchant_reject"
	AdminActionConfigChange      AdminActionType = "config_change"
	AdminActionBulkMessage       AdminActionType = "bulk_message"
	AdminActionReportGenerate    AdminActionType = "report_generate"
)

// AdminAction represents an admin action audit record
type AdminAction struct {
	ActionID   uuid.UUID       `db:"action_id" json:"action_id"`
	AdminID    uuid.UUID       `db:"admin_id" json:"admin_id"`
	ActionType AdminActionType `db:"action_type" json:"action_type"`
	TargetType string          `db:"target_type" json:"target_type"` // user, trade, withdrawal, etc.
	TargetID   uuid.UUID       `db:"target_id" json:"target_id"`
	Notes      *string         `db:"notes" json:"notes,omitempty"`
	Metadata   *string         `db:"metadata" json:"metadata,omitempty"` // JSONB
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
}

// IsUserAction checks if the action targets a user
func (a *AdminAction) IsUserAction() bool {
	return a.TargetType == "user"
}

// IsTradeAction checks if the action targets a trade
func (a *AdminAction) IsTradeAction() bool {
	return a.TargetType == "trade"
}

// AuditLog represents a generic audit log entry
type AuditLog struct {
	LogID      uuid.UUID  `db:"log_id" json:"log_id"`
	ActorID    *uuid.UUID `db:"actor_id" json:"actor_id,omitempty"` // User or admin who performed action
	Action     string     `db:"action" json:"action"`
	Resource   string     `db:"resource" json:"resource"` // table name
	ResourceID *string    `db:"resource_id" json:"resource_id,omitempty"`
	Changes    *string    `db:"changes" json:"changes,omitempty"`   // JSONB: before/after
	Metadata   *string    `db:"metadata" json:"metadata,omitempty"` // JSONB
	IPAddress  *string    `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent  *string    `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

// HasChanges checks if the audit log has recorded changes
func (a *AuditLog) HasChanges() bool {
	return a.Changes != nil && len(*a.Changes) > 0
}
