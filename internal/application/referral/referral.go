package referral

import (
	"context"

	domain "cryplio/internal/domain/referral"

	"github.com/google/uuid"
)

// TrackUseCase records a referral when a new user registers with a referral code.
type TrackUseCase struct {
	referralService domain.Service
}

func NewTrackUseCase(svc domain.Service) *TrackUseCase {
	return &TrackUseCase{referralService: svc}
}

type TrackInput struct {
	ReferrerID uuid.UUID
	RefereeID  uuid.UUID
	Code       string
}

func (uc *TrackUseCase) Execute(ctx context.Context, input TrackInput) (*domain.Referral, error) {
	return uc.referralService.Track(ctx, input.ReferrerID, input.RefereeID, input.Code)
}

// PayoutUseCase marks a referral reward as paid after a qualifying trade completes.
type PayoutUseCase struct {
	referralService domain.Service
}

func NewPayoutUseCase(svc domain.Service) *PayoutUseCase {
	return &PayoutUseCase{referralService: svc}
}

type PayoutInput struct {
	ReferralID uuid.UUID
}

func (uc *PayoutUseCase) Execute(ctx context.Context, input PayoutInput) error {
	return uc.referralService.MarkPaid(ctx, input.ReferralID)
}
