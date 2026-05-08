package wallet

import (
	"context"

	domain "cryplio/internal/domain/wallet"

	"github.com/google/uuid"
)

// DepositUseCase fetches the deposit address for a user's wallet.
type DepositUseCase struct {
	walletService domain.Service
}

func NewDepositUseCase(svc domain.Service) *DepositUseCase {
	return &DepositUseCase{walletService: svc}
}

type DepositInput struct {
	UserID       uuid.UUID
	CryptoSymbol string
}

func (uc *DepositUseCase) Execute(ctx context.Context, input DepositInput) (*domain.Wallet, error) {
	return uc.walletService.GetDepositAddress(ctx, input.UserID, input.CryptoSymbol)
}

// WithdrawUseCase initiates a withdrawal from a user's wallet.
type WithdrawUseCase struct {
	walletService domain.Service
}

func NewWithdrawUseCase(svc domain.Service) *WithdrawUseCase {
	return &WithdrawUseCase{walletService: svc}
}

type WithdrawInput struct {
	UserID       uuid.UUID
	CryptoSymbol string
	Destination  string
	Amount       float64
	Fee          float64
	Memo         *string
}

func (uc *WithdrawUseCase) Execute(ctx context.Context, input WithdrawInput) (*domain.WalletTransaction, error) {
	return uc.walletService.Withdraw(ctx, input.UserID, input.CryptoSymbol, input.Destination, input.Amount, input.Fee, input.Memo)
}

// BalanceUseCase fetches all wallet balances for a user.
type BalanceUseCase struct {
	walletService domain.Service
}

func NewBalanceUseCase(svc domain.Service) *BalanceUseCase {
	return &BalanceUseCase{walletService: svc}
}

type BalanceInput struct {
	UserID uuid.UUID
}

func (uc *BalanceUseCase) Execute(ctx context.Context, input BalanceInput) ([]domain.Wallet, error) {
	return uc.walletService.GetBalances(ctx, input.UserID)
}
