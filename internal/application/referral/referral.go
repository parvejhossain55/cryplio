package referral

import (
	"context"
	"errors"
)

// TrackUseCase coordinates referral attribution.
type TrackUseCase struct{}

func NewTrackUseCase() *TrackUseCase {
	return &TrackUseCase{}
}

type TrackInput struct {
	UserID       string
	ReferralCode string
}

func (uc *TrackUseCase) Execute(context.Context, TrackInput) error {
	return errors.New("referral track use case not implemented")
}

// PayoutUseCase coordinates referral payouts.
type PayoutUseCase struct{}

func NewPayoutUseCase() *PayoutUseCase {
	return &PayoutUseCase{}
}

type PayoutInput struct {
	UserID string
}

func (uc *PayoutUseCase) Execute(context.Context, PayoutInput) error {
	return errors.New("referral payout use case not implemented")
}
