package trade

import (
	"context"
	"errors"

	trading "cryplio/internal/domain/trading"
)

// CreateAdUseCase coordinates trade ad creation.
type CreateAdUseCase struct{}

func NewCreateAdUseCase() *CreateAdUseCase {
	return &CreateAdUseCase{}
}

type CreateAdInput struct {
	Ad *trading.TradeAd
}

func (uc *CreateAdUseCase) Execute(context.Context, CreateAdInput) (*trading.TradeAd, error) {
	return nil, errors.New("create ad use case not implemented")
}
