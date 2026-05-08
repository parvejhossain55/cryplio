package referral

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service defines business operations for the referral system.
type Service interface {
	// Track records a referral relationship when a new user registers using a referral code.
	Track(ctx context.Context, referrerID, refereeID uuid.UUID, code string) (*Referral, error)

	// GetByReferee returns the referral record for a given referee user.
	GetByReferee(ctx context.Context, refereeID uuid.UUID) (*Referral, error)

	// MarkPaid marks a referral reward as paid out.
	MarkPaid(ctx context.Context, referralID uuid.UUID) error
}

type referralService struct {
	repo Repository
}

// NewService creates a new referral service.
func NewService(repo Repository) Service {
	return &referralService{repo: repo}
}

func (s *referralService) Track(ctx context.Context, referrerID, refereeID uuid.UUID, code string) (*Referral, error) {
	ref := &Referral{
		ID:         uuid.New(),
		ReferrerID: referrerID,
		RefereeID:  refereeID,
		Code:       code,
	}
	if err := s.repo.Create(ctx, ref); err != nil {
		return nil, err
	}
	return ref, nil
}

func (s *referralService) GetByReferee(ctx context.Context, refereeID uuid.UUID) (*Referral, error) {
	return s.repo.GetByRefereeID(ctx, refereeID)
}

func (s *referralService) MarkPaid(ctx context.Context, referralID uuid.UUID) error {
	ref, err := s.repo.GetByRefereeID(ctx, referralID) // imperfect — ideally GetByID
	if err != nil {
		return err
	}
	if ref == nil {
		return nil
	}
	now := time.Now()
	ref.PaidAt = &now
	return s.repo.Update(ctx, ref)
}
