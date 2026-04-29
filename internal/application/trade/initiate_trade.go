package trade

import (
	"context"
	"errors"

	trading "cryplio/internal/domain/trading"
)

// InitiateTradeUseCase coordinates trade creation from an ad.
type InitiateTradeUseCase struct{}

func NewInitiateTradeUseCase() *InitiateTradeUseCase {
	return &InitiateTradeUseCase{}
}

type InitiateTradeInput struct {
	Trade *trading.Trade
}

func (uc *InitiateTradeUseCase) Execute(context.Context, InitiateTradeInput) (*trading.Trade, error) {
	return nil, errors.New("initiate trade use case not implemented")
}
