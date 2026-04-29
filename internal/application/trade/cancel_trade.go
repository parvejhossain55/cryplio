package trade

import (
	"context"
	"errors"

	trading "cryplio/internal/domain/trading"
	"github.com/google/uuid"
)

// CancelTradeUseCase coordinates cancellation flows.
type CancelTradeUseCase struct{}

func NewCancelTradeUseCase() *CancelTradeUseCase {
	return &CancelTradeUseCase{}
}

type CancelTradeInput struct {
	TradeID uuid.UUID
	Reason  string
}

func (uc *CancelTradeUseCase) Execute(context.Context, CancelTradeInput) (*trading.Trade, error) {
	return nil, errors.New("cancel trade use case not implemented")
}
