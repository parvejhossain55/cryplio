package trade

import (
	"context"
	"errors"

	trading "cryplio/internal/domain/trading"
	"github.com/google/uuid"
)

// ReleaseEscrowUseCase coordinates escrow release after payment confirmation.
type ReleaseEscrowUseCase struct{}

func NewReleaseEscrowUseCase() *ReleaseEscrowUseCase {
	return &ReleaseEscrowUseCase{}
}

type ReleaseEscrowInput struct {
	TradeID uuid.UUID
}

func (uc *ReleaseEscrowUseCase) Execute(context.Context, ReleaseEscrowInput) (*trading.Trade, error) {
	return nil, errors.New("release escrow use case not implemented")
}
