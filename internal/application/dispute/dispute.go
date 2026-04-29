package dispute

import (
	"context"
	"errors"

	domain "cryplio/internal/domain/dispute"
	"github.com/google/uuid"
)

// RaiseUseCase coordinates dispute creation.
type RaiseUseCase struct{}

func NewRaiseUseCase() *RaiseUseCase {
	return &RaiseUseCase{}
}

type RaiseInput struct {
	TradeID uuid.UUID
	UserID  uuid.UUID
}

func (uc *RaiseUseCase) Execute(context.Context, RaiseInput) (*domain.Dispute, error) {
	return nil, errors.New("raise dispute use case not implemented")
}

// AssignUseCase coordinates admin assignment.
type AssignUseCase struct{}

func NewAssignUseCase() *AssignUseCase {
	return &AssignUseCase{}
}

type AssignInput struct {
	DisputeID uuid.UUID
	AdminID   uuid.UUID
}

func (uc *AssignUseCase) Execute(context.Context, AssignInput) (*domain.Dispute, error) {
	return nil, errors.New("assign dispute use case not implemented")
}

// ResolveUseCase coordinates resolution decisions.
type ResolveUseCase struct{}

func NewResolveUseCase() *ResolveUseCase {
	return &ResolveUseCase{}
}

type ResolveInput struct {
	DisputeID  uuid.UUID
	AdminID    uuid.UUID
	Resolution domain.DisputeResolution
}

func (uc *ResolveUseCase) Execute(context.Context, ResolveInput) (*domain.Dispute, error) {
	return nil, errors.New("resolve dispute use case not implemented")
}
