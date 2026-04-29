package user

import (
	"context"

	identity "cryplio/internal/domain/identity"
)

// RegisterUserUseCase coordinates user registration.
type RegisterUserUseCase struct {
	authService identity.AuthService
}

func NewRegisterUserUseCase(authService identity.AuthService) *RegisterUserUseCase {
	return &RegisterUserUseCase{authService: authService}
}

type RegisterUserInput struct {
	Email    string
	Username string
	Password string
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*identity.User, error) {
	return uc.authService.Register(ctx, input.Email, input.Username, input.Password)
}
