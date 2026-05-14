package trading

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

// TradeAd represents a trade advertisement — aligned with the live trade_ads table.
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
	PaymentMethods       pq.Int64Array `db:"payment_methods" json:"payment_methods"`
	PaymentWindowMinutes int           `db:"payment_window_minutes" json:"payment_window_minutes"`
	TradeTerms           *string       `db:"trade_terms" json:"trade_terms,omitempty"`
	IsPublic             bool          `db:"is_public" json:"is_public"`
	IsPaused             bool          `db:"is_paused" json:"is_paused"`
	Status               TradeAdStatus `db:"status" json:"status"`
	CreatedAt            time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time     `db:"updated_at" json:"updated_at"`

	// Enriched fields from joins (not persisted in trade_ads table)
	Username           string         `json:"username,omitempty"`
	UserAvatar         string         `json:"user_avatar,omitempty"`
	UserLastSeen       *time.Time     `json:"user_last_seen,omitempty"`
	UserTrades         int            `json:"user_trades,omitempty"`
	UserRating         float64        `json:"user_rating,omitempty"`
	CryptoSymbol       string         `json:"crypto_symbol,omitempty"`
	FiatSymbol         string         `json:"fiat_symbol,omitempty"`
	PaymentMethodNames pq.StringArray `json:"payment_method_names,omitempty"`
}

// IsActive checks if the ad is active
func (a *TradeAd) IsActive() bool {
	return a.Status == TradeAdStatusActive && !a.IsPaused
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

// Trade represents a trade execution — aligned with the live trades table.
type Trade struct {
	TradeID              uuid.UUID       `db:"trade_id" json:"trade_id"`
	AdID                 uuid.UUID       `db:"ad_id" json:"ad_id"`
	BuyerID              uuid.UUID       `db:"buyer_id" json:"buyer_id"`
	SellerID             uuid.UUID       `db:"seller_id" json:"seller_id"`
	CryptoID             int             `db:"crypto_id" json:"crypto_id"`
	FiatID               int             `db:"fiat_id" json:"fiat_id"`
	CryptoAmount         float64         `db:"crypto_amount" json:"crypto_amount"`
	FiatAmount           float64         `db:"fiat_amount" json:"fiat_amount"`
	ExchangeRate         float64         `db:"exchange_rate" json:"exchange_rate,rate"`
	PaymentMethod        int             `db:"payment_method" json:"payment_method"`
	PaymentDetails       json.RawMessage `db:"payment_details" json:"payment_details,omitempty"`
	PriceType            PriceType       `db:"price_type" json:"price_type"`
	AgreedPrice          float64         `db:"agreed_price" json:"agreed_price"`
	Status               TradeStatus     `db:"status" json:"status"`
	PaymentWindowMinutes int             `db:"payment_window_minutes" json:"payment_window_minutes"`
	StartedAt            *time.Time      `db:"started_at" json:"started_at,omitempty"`
	PaymentMarkedAt      *time.Time      `db:"payment_marked_at" json:"payment_marked_at,omitempty"`
	ReleasedAt           *time.Time      `db:"released_at" json:"released_at,omitempty"`
	CancelledAt          *time.Time      `db:"cancelled_at" json:"cancelled_at,omitempty"`
	CompletedAt          *time.Time      `db:"completed_at" json:"completed_at,omitempty"`
	ExpiredAt            *time.Time      `db:"expired_at" json:"expired_at,omitempty"`
	CancelReason         *string         `db:"cancel_reason" json:"cancel_reason,omitempty"`
	CreatedAt            time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time       `db:"updated_at" json:"updated_at"`

	// Enriched fields
	BuyerUsername     string     `json:"buyer_username,omitempty"`
	SellerUsername    string     `json:"seller_username,omitempty"`
	CryptoSymbol      string     `json:"crypto_symbol,omitempty"`
	FiatSymbol        string     `json:"fiat_symbol,omitempty"`
	PaymentMethodName string     `json:"payment_method_name,omitempty"`
	TimerExpiresAt    *time.Time `json:"timer_expires_at,omitempty"`
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

// TradeMessage represents a chat message in a trade
type TradeMessage struct {
	ID        uuid.UUID `db:"id" json:"id"`
	TradeID   uuid.UUID `db:"trade_id" json:"trade_id"`
	SenderID  uuid.UUID `db:"sender_id" json:"sender_id"`
	Message   string    `db:"message" json:"message"`
	IsSystem  bool      `db:"is_system" json:"is_system"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
