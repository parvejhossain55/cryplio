package user

import (
	"context"

	identity "cryplio/internal/domain/identity"
)

// LoginUseCase coordinates user authentication.
type LoginUseCase struct {
	authenticator identity.Authenticator
}

func NewLoginUseCase(authenticator identity.Authenticator) *LoginUseCase {
	return &LoginUseCase{authenticator: authenticator}
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
	access, _, user, err := uc.authenticator.Login(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{Token: access, User: user}, nil
}
