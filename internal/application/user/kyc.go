package user

import (
	"context"
	"errors"

	kyc "cryplio/internal/domain/kyc"
	"github.com/google/uuid"
)

// SubmitKYCUseCase coordinates KYC submission workflows.
type SubmitKYCUseCase struct{}

func NewSubmitKYCUseCase() *SubmitKYCUseCase {
	return &SubmitKYCUseCase{}
}

type SubmitKYCInput struct {
	UserID      uuid.UUID
	DocumentURL string
}

func (uc *SubmitKYCUseCase) Execute(context.Context, SubmitKYCInput) (*kyc.KYCRecord, error) {
	return nil, errors.New("submit kyc use case not implemented")
}

// VerifyKYCUseCase coordinates KYC review and provider verification workflows.
type VerifyKYCUseCase struct{}

func NewVerifyKYCUseCase() *VerifyKYCUseCase {
	return &VerifyKYCUseCase{}
}

type VerifyKYCInput struct {
	KYCID uuid.UUID
}

func (uc *VerifyKYCUseCase) Execute(context.Context, VerifyKYCInput) (*kyc.KYCRecord, error) {
	return nil, errors.New("verify kyc use case not implemented")
}
