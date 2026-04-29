package trade

import (
	"context"
	"errors"

	trading "cryplio/internal/domain/trading"
	"github.com/google/uuid"
)

// MarkPaidUseCase coordinates payment-marked transitions.
type MarkPaidUseCase struct{}

func NewMarkPaidUseCase() *MarkPaidUseCase {
	return &MarkPaidUseCase{}
}

type MarkPaidInput struct {
	TradeID uuid.UUID
}

func (uc *MarkPaidUseCase) Execute(context.Context, MarkPaidInput) (*trading.Trade, error) {
	return nil, errors.New("mark paid use case not implemented")
}
