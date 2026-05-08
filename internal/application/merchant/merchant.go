// Package merchant contains use cases for merchant-specific operations.
// Merchant functionality (applications, verification, analytics) is planned
// for a future milestone. The use cases below are scaffolded and ready for
// implementation once the merchant domain service is available.
package merchant

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned by merchant use cases that are not yet implemented.
var ErrNotImplemented = errors.New("merchant feature not yet implemented")

// ApplyUseCase handles a user applying to become a verified merchant.
type ApplyUseCase struct{}

func NewApplyUseCase() *ApplyUseCase { return &ApplyUseCase{} }

type ApplyInput struct {
	UserID       string
	BusinessName string
	// TODO: add KYC fields, business documents, etc.
}

func (uc *ApplyUseCase) Execute(ctx context.Context, input ApplyInput) error {
	// TODO: implement when MerchantService is available
	return ErrNotImplemented
}

// VerifyUseCase handles admin verification of a merchant application.
type VerifyUseCase struct{}

func NewVerifyUseCase() *VerifyUseCase { return &VerifyUseCase{} }

type VerifyInput struct {
	UserID  string
	Approve bool
	Reason  string
}

func (uc *VerifyUseCase) Execute(ctx context.Context, input VerifyInput) error {
	// TODO: implement when MerchantService is available
	return ErrNotImplemented
}

// DashboardUseCase fetches merchant dashboard statistics.
type DashboardUseCase struct{}

func NewDashboardUseCase() *DashboardUseCase { return &DashboardUseCase{} }

type DashboardInput struct {
	UserID string
}

func (uc *DashboardUseCase) Execute(ctx context.Context, input DashboardInput) error {
	// TODO: implement when MerchantService is available
	return ErrNotImplemented
}
