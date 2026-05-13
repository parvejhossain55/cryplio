package admin

import (
	"time"

	"github.com/google/uuid"
)

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
