package notification

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents notification categories
type NotificationType string

const (
	NotificationTypeTradeStarted        NotificationType = "trade_started"
	NotificationTypeTradePaid           NotificationType = "trade_paid"
	NotificationTypeTradeReleased       NotificationType = "trade_released"
	NotificationTypeTradeCancelled      NotificationType = "trade_cancelled"
	NotificationTypeTradeDisputed       NotificationType = "trade_disputed"
	NotificationTypeDisputeResolved     NotificationType = "dispute_resolved"
	NotificationTypeNewMessage          NotificationType = "new_message"
	NotificationTypeDepositReceived     NotificationType = "deposit_received"
	NotificationTypeWithdrawalCompleted NotificationType = "withdrawal_completed"
	NotificationTypeKYCApproved         NotificationType = "kyc_approved"
	NotificationTypeKYCRejected         NotificationType = "kyc_rejected"
	NotificationTypeMerchantApproved    NotificationType = "merchant_approved"
	NotificationTypeReferralEarned      NotificationType = "referral_earned"
	NotificationTypeSystemAnnouncement  NotificationType = "system_announcement"
)

// Notification represents an in-app notification
type Notification struct {
	NotificationID uuid.UUID        `db:"notification_id" json:"notification_id"`
	UserID         uuid.UUID        `db:"user_id" json:"user_id"`
	Type           NotificationType `db:"type" json:"type"`
	Title          string           `db:"title" json:"title"`
	Message        string           `db:"message" json:"message"`
	Data           *string          `db:"data" json:"data,omitempty"` // JSON payload
	IsRead         bool             `db:"is_read" json:"is_read"`
	ReadAt         *time.Time       `db:"read_at" json:"read_at,omitempty"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	if !n.IsRead {
		n.IsRead = true
		now := time.Now()
		n.ReadAt = &now
	}
}

// IsUnread checks if the notification is unread
func (n *Notification) IsUnread() bool {
	return !n.IsRead
}

// IsRecent checks if the notification is recent (within 24 hours)
func (n *Notification) IsRecent() bool {
	return time.Since(n.CreatedAt) < 24*time.Hour
}

// NotificationPreference holds user notification settings
type NotificationPreference struct {
	UserID    uuid.UUID                 `db:"user_id" json:"user_id"`
	Email     map[NotificationType]bool `db:"-" json:"email,omitempty"` // Stored as JSONB
	Push      map[NotificationType]bool `db:"-" json:"push,omitempty"`  // Stored as JSONB
	SMS       map[NotificationType]bool `db:"-" json:"sms,omitempty"`   // Stored as JSONB
	CreatedAt time.Time                 `db:"created_at" json:"created_at"`
	UpdatedAt time.Time                 `db:"updated_at" json:"updated_at"`
}

// IsEnabledForChannel checks if a notification type is enabled for a channel
func (p *NotificationPreference) IsEnabledForChannel(nt NotificationType, channel string) bool {
	switch channel {
	case "email":
		return p.Email[nt]
	case "push":
		return p.Push[nt]
	case "sms":
		return p.SMS[nt]
	default:
		return false
	}
}

// EnableChannel enables a notification type for a channel
func (p *NotificationPreference) EnableChannel(nt NotificationType, channel string) {
	switch channel {
	case "email":
		if p.Email == nil {
			p.Email = make(map[NotificationType]bool)
		}
		p.Email[nt] = true
	case "push":
		if p.Push == nil {
			p.Push = make(map[NotificationType]bool)
		}
		p.Push[nt] = true
	case "sms":
		if p.SMS == nil {
			p.SMS = make(map[NotificationType]bool)
		}
		p.SMS[nt] = true
	}
}

// DisableChannel disables a notification type for a channel
func (p *NotificationPreference) DisableChannel(nt NotificationType, channel string) {
	switch channel {
	case "email":
		if p.Email != nil {
			p.Email[nt] = false
		}
	case "push":
		if p.Push != nil {
			p.Push[nt] = false
		}
	case "sms":
		if p.SMS != nil {
			p.SMS[nt] = false
		}
	}
}
