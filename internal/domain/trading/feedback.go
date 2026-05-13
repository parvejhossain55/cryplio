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

// TradeFeedback represents user feedback for a trade
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
