package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	identity "cryplio/internal/domain/identity"
	kyc "cryplio/internal/domain/kyc"
	kycinfra "cryplio/internal/infrastructure/kyc"

	"github.com/google/uuid"
)

// SubmitKYCUseCase coordinates KYC submission workflows.
type SubmitKYCUseCase struct {
	kycRepo    kyc.KYCRepository
	userRepo   identity.UserRepository
	persona    kycinfra.PersonaClient
	templateID string
}

func NewSubmitKYCUseCase(
	kycRepo kyc.KYCRepository,
	userRepo identity.UserRepository,
	persona kycinfra.PersonaClient,
	templateID string,
) *SubmitKYCUseCase {
	return &SubmitKYCUseCase{
		kycRepo:    kycRepo,
		userRepo:   userRepo,
		persona:    persona,
		templateID: templateID,
	}
}

type SubmitKYCInput struct {
	UserID           uuid.UUID
	Level            identity.KYCLevel
	DocumentType     string
	DocumentFrontURL string
	DocumentBackURL  *string
	SelfieURL        string
	Provider         string
}

func (uc *SubmitKYCUseCase) Execute(ctx context.Context, input SubmitKYCInput) (*kyc.KYCRecord, error) {
	// 1. Verify user exists
	user, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 2. Check if there's already a pending or approved record for this level
	existing, err := uc.kycRepo.GetLatestByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("check existing kyc: %w", err)
	}
	if existing != nil {
		if existing.IsPending() {
			return nil, errors.New("a kyc verification is already pending")
		}
		if existing.IsApproved() && existing.Level == input.Level {
			return nil, errors.New("kyc level already approved")
		}
	}

	// 3. Create KYC record
	record := &kyc.KYCRecord{
		KYCID:            uuid.New(),
		UserID:           input.UserID,
		Level:            input.Level,
		Status:           kyc.KYCStatusPending,
		DocumentType:     input.DocumentType,
		DocumentFrontURL: input.DocumentFrontURL,
		DocumentBackURL:  input.DocumentBackURL,
		SelfieURL:        input.SelfieURL,
		Provider:         input.Provider,
		SubmittedAt:      time.Now(),
	}

	// 4. If provider is persona, submit to persona
	if input.Provider == "persona" {
		externalID, err := uc.persona.CreateInquiry(ctx, input.UserID.String(), uc.templateID)
		if err != nil {
			return nil, fmt.Errorf("persona inquiry creation: %w", err)
		}
		record.ProviderReference = &externalID
	}

	// 5. Save to database
	if err := uc.kycRepo.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("save kyc record: %w", err)
	}

	return record, nil
}

// VerifyKYCUseCase coordinates KYC review and provider verification workflows.
type VerifyKYCUseCase struct {
	kycRepo  kyc.KYCRepository
	userRepo identity.UserRepository
}

func NewVerifyKYCUseCase(
	kycRepo kyc.KYCRepository,
	userRepo identity.UserRepository,
) *VerifyKYCUseCase {
	return &VerifyKYCUseCase{
		kycRepo:  kycRepo,
		userRepo: userRepo,
	}
}

type VerifyKYCInput struct {
	KYCID           uuid.UUID
	AdminID         uuid.UUID
	Approved        bool
	RejectionReason string
}

func (uc *VerifyKYCUseCase) Execute(ctx context.Context, input VerifyKYCInput) (*kyc.KYCRecord, error) {
	// 1. Get KYC record
	record, err := uc.kycRepo.GetByID(ctx, input.KYCID)
	if err != nil {
		return nil, fmt.Errorf("get kyc record: %w", err)
	}
	if record == nil {
		return nil, errors.New("kyc record not found")
	}

	if !record.IsPending() {
		return nil, errors.New("kyc record is not in pending status")
	}

	// 2. Process approval/rejection
	if input.Approved {
		record.Approve(input.AdminID)

		// Upgrade user KYC level
		user, err := uc.userRepo.GetByID(ctx, record.UserID)
		if err != nil {
			return nil, fmt.Errorf("get user: %w", err)
		}
		if user != nil {
			user.UpgradeKYC(record.Level)
			if err := uc.userRepo.Update(ctx, user); err != nil {
				return nil, fmt.Errorf("update user kyc level: %w", err)
			}
		}

		// TODO: Automated AML check (FR-128)
		// For now just mark as screened
		record.SetAMLResult(map[string]string{"status": "clear", "source": "internal_mock"})
	} else {
		record.Reject(input.AdminID, input.RejectionReason)
	}

	// 3. Update record
	if err := uc.kycRepo.Update(ctx, record); err != nil {
		return nil, fmt.Errorf("update kyc record: %w", err)
	}

	return record, nil
}
