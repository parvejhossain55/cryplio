package escrow

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, lock *Lock) error
	Update(ctx context.Context, lock *Lock) error
	GetByTradeID(ctx context.Context, tradeID uuid.UUID) (*Lock, error)
}
