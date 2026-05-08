package identity

import (
	"context"
	"fmt"

	"cryplio/pkg/apperrors"

	"github.com/google/uuid"
)

// ─── User Payment Methods ─────────────────────────────────────────────────────
//
// These methods manage payment method profiles attached to a user account
// (e.g. bKash number, bank account details). They are distinct from the
// platform-level payment method catalogue managed by admins.

// AddPaymentMethod stores a new payment method for the user.
func (s *authService) AddPaymentMethod(ctx context.Context, userID uuid.UUID, pm *UserPaymentMethod) (*UserPaymentMethod, error) {
	pm.UserID = userID
	if err := s.userRepo.CreateUserPaymentMethod(ctx, pm); err != nil {
		return nil, fmt.Errorf("create payment method: %w", err)
	}
	return pm, nil
}

// GetPaymentMethods returns all payment methods saved by the user.
func (s *authService) GetPaymentMethods(ctx context.Context, userID uuid.UUID) ([]UserPaymentMethod, error) {
	return s.userRepo.GetUserPaymentMethods(ctx, userID)
}

// UpdatePaymentMethod replaces the fields of an existing payment method.
func (s *authService) UpdatePaymentMethod(ctx context.Context, userID uuid.UUID, pm *UserPaymentMethod) (*UserPaymentMethod, error) {
	pm.UserID = userID
	if err := s.userRepo.UpdateUserPaymentMethod(ctx, pm); err != nil {
		return nil, fmt.Errorf("update payment method: %w", err)
	}
	return pm, nil
}

// RemovePaymentMethod deletes a payment method after verifying ownership.
func (s *authService) RemovePaymentMethod(ctx context.Context, userID, pmID uuid.UUID) error {
	pm, err := s.userRepo.GetUserPaymentMethod(ctx, pmID)
	if err != nil {
		return err
	}
	if pm == nil || pm.UserID != userID {
		return apperrors.NotFound("payment method not found", nil)
	}
	return s.userRepo.DeleteUserPaymentMethod(ctx, pmID)
}

// SetDefaultPaymentMethod marks one payment method as the user's default.
func (s *authService) SetDefaultPaymentMethod(ctx context.Context, userID, pmID uuid.UUID) error {
	pm, err := s.userRepo.GetUserPaymentMethod(ctx, pmID)
	if err != nil {
		return err
	}
	if pm == nil || pm.UserID != userID {
		return apperrors.NotFound("payment method not found", nil)
	}
	return s.userRepo.SetDefaultUserPaymentMethod(ctx, userID, pmID)
}
