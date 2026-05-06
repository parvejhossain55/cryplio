package trading

import (
	"time"

	"github.com/google/uuid"
)

// AdType represents trade ad type
type AdType string

const (
	AdTypeBuy  AdType = "buy"
	AdTypeSell AdType = "sell"
)

// PriceType represents price specification
type PriceType string

const (
	PriceTypeFixed    PriceType = "fixed"
	PriceTypeFloating PriceType = "floating"
)

// TradeAdStatus represents the status of a trade advertisement.
type TradeAdStatus string

const (
	TradeAdStatusActive TradeAdStatus = "active"
	TradeAdStatusPaused TradeAdStatus = "paused"
	TradeAdStatusDraft  TradeAdStatus = "draft"
	TradeAdStatusClosed TradeAdStatus = "closed"
)

// TradeAd represents a trade advertisement
type TradeAd struct {
	AdID                 uuid.UUID     `db:"ad_id" json:"ad_id"`
	UserID               uuid.UUID     `db:"user_id" json:"user_id"`
	Type                 AdType        `db:"type" json:"type"`
	CryptoID             int           `db:"crypto_id" json:"crypto_id"`
	FiatID               int           `db:"fiat_id" json:"fiat_id"`
	PriceType            PriceType     `db:"price_type" json:"price_type"`
	Price                float64       `db:"price" json:"price"`
	FloatingMarkup       *float64      `db:"floating_markup" json:"floating_markup,omitempty"`
	MinAmount            float64       `db:"min_amount" json:"min_amount"`
	MaxAmount            float64       `db:"max_amount" json:"max_amount"`
	PaymentMethods       []int         `db:"payment_methods" json:"payment_methods"` // PostgreSQL int[]
	TradeTerms           *string       `db:"trade_terms" json:"trade_terms,omitempty"`
	PaymentWindowMinutes int           `db:"payment_window_minutes" json:"payment_window_minutes"`
	IsPublic             bool          `db:"is_public" json:"is_public"`
	IsPaused             bool          `db:"is_paused" json:"is_paused"`
	VisibilityStartAt    *time.Time    `db:"visibility_start_at" json:"visibility_start_at,omitempty"`
	VisibilityEndAt      *time.Time    `db:"visibility_end_at" json:"visibility_end_at,omitempty"`
	Timezone             string        `db:"timezone" json:"timezone"`
	AutoRepost           bool          `db:"auto_repost" json:"auto_repost"`
	RepostCount          int           `db:"repost_count" json:"repost_count"`
	ViewsCount           int           `db:"views_count" json:"views_count"`
	ResponseCount        int           `db:"response_count" json:"response_count"`
	LockedBalance        float64       `db:"locked_balance" json:"locked_balance"`
	Status               TradeAdStatus `db:"status" json:"status"`
	FirstPublishedAt     *time.Time    `db:"first_published_at" json:"first_published_at,omitempty"`
	PublishedAt          time.Time     `db:"published_at" json:"published_at"`
	ExpiresAt            *time.Time    `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt            time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time     `db:"updated_at" json:"updated_at"`
	DeletedAt            *time.Time    `db:"deleted_at" json:"deleted_at,omitempty"`

	// Enriched fields from joins (not persisted in trade_ads table)
	Username     string     `json:"username,omitempty"`
	UserAvatar   string     `json:"user_avatar,omitempty"`
	UserLastSeen *time.Time `json:"user_last_seen,omitempty"`
	UserTrades   int        `json:"user_trades,omitempty"`
	UserRating   float64    `json:"user_rating,omitempty"`
}

// IsActive checks if the ad is active and visible
func (a *TradeAd) IsActive() bool {
	if a.IsPaused || a.DeletedAt != nil {
		return false
	}
	now := time.Now()
	if a.VisibilityStartAt != nil && now.Before(*a.VisibilityStartAt) {
		return false
	}
	if a.VisibilityEndAt != nil && now.After(*a.VisibilityEndAt) {
		return false
	}
	if a.ExpiresAt != nil && now.After(*a.ExpiresAt) {
		return false
	}
	return a.Status == TradeAdStatusActive
}

// Pause pauses the ad
func (a *TradeAd) Pause() {
	a.IsPaused = true
}

// Resume resumes the ad
func (a *TradeAd) Resume() {
	a.IsPaused = false
}

// CanAcceptAmount checks if the ad can accept the given amount
func (a *TradeAd) CanAcceptAmount(amount float64) bool {
	return amount >= a.MinAmount && amount <= a.MaxAmount
}

// TradeStatus represents trade lifecycle state
type TradeStatus string

const (
	TradeStatusPending   TradeStatus = "pending"
	TradeStatusActive    TradeStatus = "active"
	TradeStatusPaid      TradeStatus = "paid"
	TradeStatusReleased  TradeStatus = "released"
	TradeStatusCancelled TradeStatus = "cancelled"
	TradeStatusDisputed  TradeStatus = "disputed"
	TradeStatusCompleted TradeStatus = "completed"
	TradeStatusExpired   TradeStatus = "expired"
)

// Trade represents a trade execution
type Trade struct {
	TradeID                uuid.UUID   `db:"trade_id" json:"trade_id"`
	AdID                   uuid.UUID   `db:"ad_id" json:"ad_id"`
	BuyerID                uuid.UUID   `db:"buyer_id" json:"buyer_id"`
	SellerID               uuid.UUID   `db:"seller_id" json:"seller_id"`
	CryptoID               int         `db:"crypto_id" json:"crypto_id"`
	FiatID                 int         `db:"fiat_id" json:"fiat_id"`
	CryptoAmount           float64     `db:"crypto_amount" json:"crypto_amount"`
	FiatAmount             float64     `db:"fiat_amount" json:"fiat_amount"`
	ExchangeRate           float64     `db:"exchange_rate" json:"exchange_rate"`
	PaymentMethod          int         `db:"payment_method" json:"payment_method"`
	PriceType              PriceType   `db:"price_type" json:"price_type"`
	AgreedPrice            float64     `db:"agreed_price" json:"agreed_price"`
	Status                 TradeStatus `db:"status" json:"status"`
	DisputeID              *uuid.UUID  `db:"dispute_id" json:"dispute_id,omitempty"`
	ChatRoomID             *string     `db:"chat_room_id" json:"chat_room_id,omitempty"`
	StartedAt              *time.Time  `db:"started_at" json:"started_at,omitempty"`
	PaymentMarkedAt        *time.Time  `db:"payment_marked_at" json:"payment_marked_at,omitempty"`
	ReleasedAt             *time.Time  `db:"released_at" json:"released_at,omitempty"`
	CancelledAt            *time.Time  `db:"cancelled_at" json:"cancelled_at,omitempty"`
	CompletedAt            *time.Time  `db:"completed_at" json:"completed_at,omitempty"`
	ExpiredAt              *time.Time  `db:"expired_at" json:"expired_at,omitempty"`
	PaymentWindowMinutes   int         `db:"payment_window_minutes" json:"payment_window_minutes"`
	TimeRemainingSeconds   *int        `db:"time_remaining_seconds" json:"time_remaining_seconds,omitempty"`
	IsAutoDisputeTriggered bool        `db:"is_auto_dispute_triggered" json:"is_auto_dispute_triggered"`
	CancelReason           *string     `db:"cancel_reason" json:"cancel_reason,omitempty"`
	EscrowTxnHash          *string     `db:"escrow_txn_hash" json:"escrow_txn_hash,omitempty"`
	EscrowContractAddress  *string     `db:"escrow_contract_address" json:"escrow_contract_address,omitempty"`
	CreatedAt              time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt              *time.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
}

// IsActive checks if the trade is in an active state
func (t *Trade) IsActive() bool {
	return t.Status == TradeStatusPending || t.Status == TradeStatusActive || t.Status == TradeStatusPaid
}

// CanRelease checks if the trade can be released
func (t *Trade) CanRelease() bool {
	return t.Status == TradeStatusPaid && t.PaymentMarkedAt != nil
}

// CanCancel checks if the trade can be cancelled
func (t *Trade) CanCancel() bool {
	return t.Status == TradeStatusPending || t.Status == TradeStatusActive
}

// MarkAsPaid marks the trade as paid
func (t *Trade) MarkAsPaid() {
	now := time.Now()
	t.Status = TradeStatusPaid
	t.PaymentMarkedAt = &now
}

// Release the escrow to the buyer
func (t *Trade) Release() {
	now := time.Now()
	t.Status = TradeStatusReleased
	t.ReleasedAt = &now
}

// Cancel the trade
func (t *Trade) Cancel(reason string) {
	now := time.Now()
	t.Status = TradeStatusCancelled
	t.CancelledAt = &now
	t.CancelReason = &reason
}

// Complete marks the trade as completed
func (t *Trade) Complete() {
	now := time.Now()
	t.Status = TradeStatusCompleted
	t.CompletedAt = &now
}

// IsDisputed checks if the trade is in dispute
func (t *Trade) IsDisputed() bool {
	return t.Status == TradeStatusDisputed
}

// IsExpired checks if the trade has expired
func (t *Trade) IsExpired() bool {
	return t.Status == TradeStatusExpired
}

// TradeMessageType represents message content type
type TradeMessageType string

const (
	TradeMessageTypeText  TradeMessageType = "text"
	TradeMessageTypeImage TradeMessageType = "image"
	TradeMessageTypeFile  TradeMessageType = "file"
)

// TradeMessage represents a chat message in a trade
type TradeMessage struct {
	MessageID    uuid.UUID        `db:"message_id" json:"message_id"`
	TradeID      uuid.UUID        `db:"trade_id" json:"trade_id"`
	SenderID     uuid.UUID        `db:"sender_id" json:"sender_id"`
	MessageType  TradeMessageType `db:"message_type" json:"message_type"`
	Content      *string          `db:"content" json:"content,omitempty"`
	FileURL      *string          `db:"file_url" json:"file_url,omitempty"`
	FileMimeType *string          `db:"file_mime_type" json:"file_mime_type,omitempty"`
	FileSize     *int             `db:"file_size" json:"file_size,omitempty"`
	IsRead       bool             `db:"is_read" json:"is_read"`
	ReadAt       *time.Time       `db:"read_at" json:"read_at,omitempty"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
	DeletedAt    *time.Time       `db:"deleted_at" json:"deleted_at,omitempty"`
}

// MarkAsRead marks the message as read
func (m *TradeMessage) MarkAsRead() {
	if !m.IsRead {
		m.IsRead = true
		now := time.Now()
		m.ReadAt = &now
	}
}

// TradeAttachment represents an attached file to a trade message
type TradeAttachment struct {
	ID         int       `db:"id" json:"id"`
	MessageID  uuid.UUID `db:"message_id" json:"message_id"`
	FileName   string    `db:"file_name" json:"file_name"`
	FileURL    string    `db:"file_url" json:"file_url"`
	MimeType   string    `db:"mime_type" json:"mime_type"`
	SizeBytes  int64     `db:"size_bytes" json:"size_bytes"`
	UploaderID uuid.UUID `db:"uploader_id" json:"uploader_id"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}

// IsImage checks if the attachment is an image
func (a *TradeAttachment) IsImage() bool {
	return a.MimeType == "image/jpeg" || a.MimeType == "image/png" || a.MimeType == "image/gif"
}
