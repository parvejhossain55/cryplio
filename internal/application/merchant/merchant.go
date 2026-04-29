package merchant

import (
	"context"
	"errors"
)

// ApplyUseCase coordinates merchant application flows.
type ApplyUseCase struct{}

func NewApplyUseCase() *ApplyUseCase {
	return &ApplyUseCase{}
}

type ApplyInput struct {
	UserID string
}

func (uc *ApplyUseCase) Execute(context.Context, ApplyInput) error {
	return errors.New("merchant apply use case not implemented")
}

// VerifyUseCase coordinates merchant verification flows.
type VerifyUseCase struct{}

func NewVerifyUseCase() *VerifyUseCase {
	return &VerifyUseCase{}
}

type VerifyInput struct {
	UserID string
}

func (uc *VerifyUseCase) Execute(context.Context, VerifyInput) error {
	return errors.New("merchant verify use case not implemented")
}

// DashboardUseCase coordinates merchant dashboard reads.
type DashboardUseCase struct{}

func NewDashboardUseCase() *DashboardUseCase {
	return &DashboardUseCase{}
}

type DashboardInput struct {
	UserID string
}

func (uc *DashboardUseCase) Execute(context.Context, DashboardInput) error {
	return errors.New("merchant dashboard use case not implemented")
}
