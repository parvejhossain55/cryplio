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
	CreateWallet(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error)
	Withdraw(ctx context.Context, userID uuid.UUID, cryptoSymbol, destination string, amount float64, fee float64, memo *string) (*WalletTransaction, error)
	GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]WalletTransaction, int, error)
	ListPendingWithdrawals(ctx context.Context, limit, offset int) ([]WalletTransaction, int, error)
	ApproveWithdrawal(ctx context.Context, txID, adminID uuid.UUID, txHash string) error
	RejectWithdrawal(ctx context.Context, txID, adminID uuid.UUID, reason string) error
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

func (s *service) CreateWallet(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error) {
	symbol := strings.TrimSpace(strings.ToUpper(cryptoSymbol))
	if symbol == "" {
		return nil, apperrors.InvalidInput("crypto symbol is required", nil)
	}

	// Get crypto_id from database (handles duplicates by picking first match)
	cryptoID, err := s.repo.GetCryptoIDBySymbol(ctx, symbol)
	if err != nil {
		return nil, apperrors.NotFound("cryptocurrency not supported", err)
	}

	// Check if wallet already exists using crypto_id
	existing, err := s.repo.GetByUserAndCryptoID(ctx, userID, cryptoID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperrors.InvalidInput("wallet already exists for this cryptocurrency", nil)
	}

	// Generate a blockchain address
	address := fmt.Sprintf("0x%s%s%d", userID.String()[:8], symbol, time.Now().Unix())

	wallet := &Wallet{
		WalletID:      uuid.New(),
		UserID:        userID,
		CryptoID:      cryptoID,
		Address:       address,
		Balance:       0,
		LockedBalance: 0,
		IsActive:      true,
		LastUpdated:   time.Now(),
		CreatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, apperrors.InvalidInput("wallet already exists for this cryptocurrency", nil)
		}
		return nil, fmt.Errorf("create wallet: %w", err)
	}

	return wallet, nil
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

	// Check daily withdrawal limit ($500 USD equivalent)
	dailyTotal, err := s.repo.GetDailyWithdrawalTotal(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("check daily withdrawal limit: %w", err)
	}

	const dailyLimitUSD = 500.0
	// For MVP, we assume USDT = USD for simplicity. In production, convert based on current rates
	if dailyTotal+amount > dailyLimitUSD {
		return nil, apperrors.InvalidInput(fmt.Sprintf("daily withdrawal limit exceeded (max $%.0f)", dailyLimitUSD), nil)
	}

	wallet, err := s.GetDepositAddress(ctx, userID, cryptoSymbol)
	if err != nil {
		return nil, err
	}

	totalDebit := amount + fee
	if err := WithdrawFunds(wallet, totalDebit); err != nil {
		return nil, apperrors.InsufficientFunds("insufficient available balance", err)
	}

	const approvalThreshold = 1000.0
	requiresApproval := amount > approvalThreshold

	now := time.Now()
	tx := &WalletTransaction{
		TxID:               uuid.New(),
		WalletID:           wallet.WalletID,
		UserID:             userID,
		Type:               TransactionTypeWithdrawal,
		Status:             TransactionStatusPending,
		Amount:             amount,
		Fee:                fee,
		NetAmount:          amount - fee,
		BalanceAfter:       wallet.Balance,
		Memo:               memo,
		RequiresApproval:   requiresApproval,
		DestinationAddress: &destination,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.Update(ctx, wallet); err != nil {
		return nil, fmt.Errorf("update wallet: %w", err)
	}

	// For withdrawals requiring approval, do not execute on-chain yet
	if !requiresApproval {
		// Execute on-chain transfer for small withdrawals
		onChainTxHash, err := s.walletClient.Send(ctx, tx, destination)
		if err != nil {
			// In a real app, we would revert the balance or mark the TX as failed
			return nil, fmt.Errorf("on-chain transfer failed: %w", err)
		}
		tx.TxHash = &onChainTxHash
		tx.Status = TransactionStatusCompleted
	}

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

func (s *service) ListPendingWithdrawals(ctx context.Context, limit, offset int) ([]WalletTransaction, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListPendingWithdrawals(ctx, limit, offset)
}

func (s *service) ApproveWithdrawal(ctx context.Context, txID, adminID uuid.UUID, txHash string) error {
	return s.repo.ApproveWithdrawal(ctx, txID, adminID, txHash)
}

func (s *service) RejectWithdrawal(ctx context.Context, txID, adminID uuid.UUID, reason string) error {
	return s.repo.RejectWithdrawal(ctx, txID, adminID, reason)
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
