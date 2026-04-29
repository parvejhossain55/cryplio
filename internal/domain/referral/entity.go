package referral

import (
	"time"

	"github.com/google/uuid"
)

type Referral struct {
	ID           uuid.UUID  `json:"id"`
	ReferrerID   uuid.UUID  `json:"referrer_id"`
	RefereeID    uuid.UUID  `json:"referee_id"`
	Code         string     `json:"code"`
	RewardAmount float64    `json:"reward_amount"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
