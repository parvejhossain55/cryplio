package dispute

import (
	"context"

	domain "cryplio/internal/domain/dispute"

	"github.com/google/uuid"
)

// RaiseUseCase orchestrates raising a new dispute on an active trade.
type RaiseUseCase struct {
	disputeService domain.Service
}

func NewRaiseUseCase(svc domain.Service) *RaiseUseCase {
	return &RaiseUseCase{disputeService: svc}
}

type RaiseInput struct {
	TradeID    uuid.UUID
	UserID     uuid.UUID
	ReasonCode string
	ReasonText string
}

func (uc *RaiseUseCase) Execute(ctx context.Context, input RaiseInput) (*domain.Dispute, error) {
	d := &domain.Dispute{
		DisputeID:  uuid.New(),
		TradeID:    input.TradeID,
		RaisedBy:   input.UserID,
		ReasonCode: input.ReasonCode,
		ReasonText: &input.ReasonText,
		Status:     domain.DisputeStatusPending,
	}
	if err := uc.disputeService.CreateDispute(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// AssignUseCase orchestrates assigning a dispute to an admin.
type AssignUseCase struct {
	disputeService domain.Service
}

func NewAssignUseCase(svc domain.Service) *AssignUseCase {
	return &AssignUseCase{disputeService: svc}
}

type AssignInput struct {
	DisputeID uuid.UUID
	AdminID   uuid.UUID
}

func (uc *AssignUseCase) Execute(ctx context.Context, input AssignInput) error {
	return uc.disputeService.AssignDispute(ctx, input.DisputeID, input.AdminID)
}

// ResolveUseCase orchestrates resolving a dispute.
type ResolveUseCase struct {
	disputeService domain.Service
}

func NewResolveUseCase(svc domain.Service) *ResolveUseCase {
	return &ResolveUseCase{disputeService: svc}
}

type ResolveInput struct {
	DisputeID  uuid.UUID
	AdminID    uuid.UUID
	Resolution domain.DisputeResolution
	Note       string
}

func (uc *ResolveUseCase) Execute(ctx context.Context, input ResolveInput) error {
	return uc.disputeService.ResolveDispute(ctx, input.DisputeID, input.AdminID, input.Resolution, input.Note)
}
