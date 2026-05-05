package wallet

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cryplio/pkg/apperrors"

	"github.com/google/uuid"
)

type Service interface {
	GetBalances(ctx context.Context, userID uuid.UUID) ([]Wallet, error)
	GetDepositAddress(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error)
	Withdraw(ctx context.Context, userID uuid.UUID, cryptoSymbol, destination string, amount float64, fee float64, memo *string) (*WalletTransaction, error)
	GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]WalletTransaction, int, error)
}

type service struct {
	repo         Repository
	walletClient WalletClient
}

func NewService(repo Repository, walletClient WalletClient) Service {
	return &service{
		repo:         repo,
		walletClient: walletClient,
	}
}

func (s *service) GetBalances(ctx context.Context, userID uuid.UUID) ([]Wallet, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) GetDepositAddress(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error) {
	symbol := strings.TrimSpace(strings.ToUpper(cryptoSymbol))
	if symbol == "" {
		return nil, apperrors.InvalidInput("crypto symbol is required", nil)
	}

	wallet, err := s.repo.GetByUserAndCrypto(ctx, userID, symbol)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, apperrors.NotFound("wallet not found for requested asset", nil)
	}

	return wallet, nil
}

func (s *service) Withdraw(ctx context.Context, userID uuid.UUID, cryptoSymbol, destination string, amount float64, fee float64, memo *string) (*WalletTransaction, error) {
	if amount <= 0 {
		return nil, apperrors.InvalidInput("amount must be greater than zero", nil)
	}
	if strings.TrimSpace(destination) == "" {
		return nil, apperrors.InvalidInput("destination is required", nil)
	}
	if fee < 0 {
		return nil, apperrors.InvalidInput("fee cannot be negative", nil)
	}

	wallet, err := s.GetDepositAddress(ctx, userID, cryptoSymbol)
	if err != nil {
		return nil, err
	}

	totalDebit := amount + fee
	if err := WithdrawFunds(wallet, totalDebit); err != nil {
		return nil, apperrors.InsufficientFunds("insufficient available balance", err)
	}

	now := time.Now()
	tx := &WalletTransaction{
		TxID:         uuid.New(),
		WalletID:     wallet.WalletID,
		UserID:       userID,
		Type:         TransactionTypeWithdrawal,
		Status:       TransactionStatusPending,
		Amount:       amount,
		Fee:          fee,
		NetAmount:    amount - fee,
		BalanceAfter: wallet.Balance,
		Memo:         memo,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Update(ctx, wallet); err != nil {
		return nil, fmt.Errorf("update wallet: %w", err)
	}

	// 8. Execute on-chain transfer
	onChainTxHash, err := s.walletClient.Send(ctx, tx, destination)
	if err != nil {
		// In a real app, we would revert the balance or mark the TX as failed
		return nil, fmt.Errorf("on-chain transfer failed: %w", err)
	}
	tx.TxHash = &onChainTxHash

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("create wallet transaction: %w", err)
	}

	return tx, nil
}

func (s *service) GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]WalletTransaction, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListTransactionsByUser(ctx, userID, limit, offset)
}

func DepositFunds(wallet *Wallet, amount float64) error {
	if wallet == nil || amount <= 0 {
		return errors.New("invalid deposit")
	}
	wallet.Credit(amount)
	return nil
}

func WithdrawFunds(wallet *Wallet, amount float64) error {
	if wallet == nil || amount <= 0 {
		return errors.New("invalid withdrawal")
	}
	return wallet.Debit(amount)
}

func LockBalance(wallet *Wallet, amount float64) error {
	if wallet == nil || amount <= 0 {
		return errors.New("invalid lock amount")
	}
	return wallet.Lock(amount)
}
