package escrow

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusLocked   Status = "locked"
	StatusReleased Status = "released"
	StatusReturned Status = "returned"
)

type Lock struct {
	ID              uuid.UUID  `json:"id"`
	TradeID         uuid.UUID  `json:"trade_id"`
	Amount          float64    `json:"amount"`
	TxHash          *string    `json:"tx_hash,omitempty"`
	ContractAddress *string    `json:"contract_address,omitempty"`
	Status          Status     `json:"status"`
	LockedAt        time.Time  `json:"locked_at"`
	ReleasedAt      *time.Time `json:"released_at,omitempty"`
	ReturnedAt      *time.Time `json:"returned_at,omitempty"`
}
