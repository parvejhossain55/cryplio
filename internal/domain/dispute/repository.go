package dispute

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, dispute *Dispute) error
	Update(ctx context.Context, dispute *Dispute) error
	GetByID(ctx context.Context, disputeID uuid.UUID) (*Dispute, error)
	List(ctx context.Context, status string, limit, offset int) ([]*Dispute, int, error)
	CountDisputes(ctx context.Context, status string) (int, error)
}
