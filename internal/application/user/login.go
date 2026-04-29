package user

import (
	"context"

	identity "cryplio/internal/domain/identity"
)

// LoginUseCase coordinates user authentication.
type LoginUseCase struct {
	authService identity.AuthService
}

func NewLoginUseCase(authService identity.AuthService) *LoginUseCase {
	return &LoginUseCase{authService: authService}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
	User  *identity.User
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	token, user, err := uc.authService.Login(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{Token: token, User: user}, nil
}
