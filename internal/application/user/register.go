package user

import (
	"context"

	identity "cryplio/internal/domain/identity"
)

// RegisterUserUseCase coordinates user registration.
type RegisterUserUseCase struct {
	registrar identity.UserRegistrar
}

func NewRegisterUserUseCase(registrar identity.UserRegistrar) *RegisterUserUseCase {
	return &RegisterUserUseCase{registrar: registrar}
}

type RegisterUserInput struct {
	Email    string
	Username string
	Password string
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*identity.User, error) {
	return uc.registrar.Register(ctx, input.Email, input.Username, input.Password)
}
