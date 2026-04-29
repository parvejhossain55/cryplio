package referral

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, referral *Referral) error
	GetByRefereeID(ctx context.Context, refereeID uuid.UUID) (*Referral, error)
	Update(ctx context.Context, referral *Referral) error
}
