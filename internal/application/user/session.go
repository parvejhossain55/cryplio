package user

import (
	"context"
	"errors"
)

// RefreshTokenUseCase coordinates access token refresh flows.
type RefreshTokenUseCase struct{}

func NewRefreshTokenUseCase() *RefreshTokenUseCase {
	return &RefreshTokenUseCase{}
}

type RefreshTokenInput struct {
	RefreshToken string
}

func (uc *RefreshTokenUseCase) Execute(context.Context, RefreshTokenInput) (string, error) {
	return "", errors.New("refresh token use case not implemented")
}

// RevokeDeviceUseCase coordinates session/device revocation flows.
type RevokeDeviceUseCase struct{}

func NewRevokeDeviceUseCase() *RevokeDeviceUseCase {
	return &RevokeDeviceUseCase{}
}

type RevokeDeviceInput struct {
	TokenID string
}

func (uc *RevokeDeviceUseCase) Execute(context.Context, RevokeDeviceInput) error {
	return errors.New("revoke device use case not implemented")
}
