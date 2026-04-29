package wallet

import (
	"context"
	"errors"

	domain "cryplio/internal/domain/wallet"
	"github.com/google/uuid"
)

// DepositUseCase coordinates deposit workflows.
type DepositUseCase struct{}

func NewDepositUseCase() *DepositUseCase {
	return &DepositUseCase{}
}

type DepositInput struct {
	WalletID uuid.UUID
	Amount   float64
}

func (uc *DepositUseCase) Execute(context.Context, DepositInput) (*domain.WalletTransaction, error) {
	return nil, errors.New("deposit use case not implemented")
}

// WithdrawUseCase coordinates withdrawal workflows.
type WithdrawUseCase struct{}

func NewWithdrawUseCase() *WithdrawUseCase {
	return &WithdrawUseCase{}
}

type WithdrawInput struct {
	WalletID    uuid.UUID
	Amount      float64
	Destination string
}

func (uc *WithdrawUseCase) Execute(context.Context, WithdrawInput) (*domain.WalletTransaction, error) {
	return nil, errors.New("withdraw use case not implemented")
}

// BalanceUseCase coordinates balance reads.
type BalanceUseCase struct{}

func NewBalanceUseCase() *BalanceUseCase {
	return &BalanceUseCase{}
}

type BalanceInput struct {
	WalletID uuid.UUID
}

func (uc *BalanceUseCase) Execute(context.Context, BalanceInput) (*domain.Wallet, error) {
	return nil, errors.New("balance use case not implemented")
}
