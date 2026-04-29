package trading

import (
	"time"

	"github.com/google/uuid"
)

// FeedbackRating represents feedback rating
type FeedbackRating string

const (
	FeedbackPositive FeedbackRating = "positive"
	FeedbackNeutral  FeedbackRating = "neutral"
	FeedbackNegative FeedbackRating = "negative"
)

// ReferralStatus represents the status of a referral reward.
type ReferralStatus string

const (
	ReferralStatusPending ReferralStatus = "pending"
	ReferralStatusPaid    ReferralStatus = "paid"
)

// TradeFeedback represents trade feedback/rating
type TradeFeedback struct {
	FeedbackID uuid.UUID      `db:"feedback_id" json:"feedback_id"`
	TradeID    uuid.UUID      `db:"trade_id" json:"trade_id"`
	FromUserID uuid.UUID      `db:"from_user_id" json:"from_user_id"`
	ToUserID   uuid.UUID      `db:"to_user_id" json:"to_user_id"`
	Rating     FeedbackRating `db:"rating" json:"rating"`
	Comment    *string        `db:"comment" json:"comment,omitempty"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
}

// IsPositive checks if the feedback is positive
func (f *TradeFeedback) IsPositive() bool {
	return f.Rating == FeedbackPositive
}

// IsNegative checks if the feedback is negative
func (f *TradeFeedback) IsNegative() bool {
	return f.Rating == FeedbackNegative
}

// IsNeutral checks if the feedback is neutral
func (f *TradeFeedback) IsNeutral() bool {
	return f.Rating == FeedbackNeutral
}

// Referral represents a user referral
type Referral struct {
	ReferralID   uuid.UUID      `db:"referral_id" json:"referral_id"`
	ReferrerID   uuid.UUID      `db:"referrer_id" json:"referrer_id"`
	RefereeID    uuid.UUID      `db:"referee_id" json:"referee_id"`
	CodeUsed     string         `db:"code_used" json:"code_used"`
	RewardAmount float64        `db:"reward_amount" json:"reward_amount"`
	RewardType   string         `db:"reward_type" json:"reward_type"` // crypto, fiat, fee_discount
	Status       ReferralStatus `db:"status" json:"status"`           // pending, paid
	PaidAt       *time.Time     `db:"paid_at" json:"paid_at,omitempty"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
}

// IsPaid checks if the referral reward has been paid
func (r *Referral) IsPaid() bool {
	return r.Status == ReferralStatusPaid && r.PaidAt != nil
}

// IsPending checks if the referral reward is pending
func (r *Referral) IsPending() bool {
	return r.Status == ReferralStatusPending
}
