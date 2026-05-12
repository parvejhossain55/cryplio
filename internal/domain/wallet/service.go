package wallet

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cryplio/internal/domain/platform"
	"cryplio/pkg/apperrors"
	"cryplio/pkg/config"
	sharedcrypto "cryplio/pkg/crypto"

	"github.com/google/uuid"
)

type Service interface {
	GetBalances(ctx context.Context, userID uuid.UUID) ([]Wallet, error)
	GetDepositAddress(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error)
	CreateDefaultWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	Withdraw(ctx context.Context, userID uuid.UUID, cryptoSymbol, destination string, amount float64, fee float64, memo *string, emailCode string) (*WalletTransaction, error)
	GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]WalletTransaction, int, error)
	ListPendingWithdrawals(ctx context.Context, limit, offset int) ([]WalletTransaction, int, error)
	ApproveWithdrawal(ctx context.Context, txID, adminID uuid.UUID, txHash string) error
	RejectWithdrawal(ctx context.Context, txID, adminID uuid.UUID, reason string) error
	GetDailyLimitInfo(ctx context.Context, userID uuid.UUID) (float64, float64, error)
}

type service struct {
	repo            Repository
	walletClient    WalletClient
	platformService platform.PlatformService
	cfg             *config.Config
}

func NewService(repo Repository, walletClient WalletClient, platformService platform.PlatformService, cfg *config.Config) Service {
	return &service{
		repo:            repo,
		walletClient:    walletClient,
		platformService: platformService,
		cfg:             cfg,
	}
}

func (s *service) GetBalances(ctx context.Context, userID uuid.UUID) ([]Wallet, error) {
	// 1. Get the single generic wallet for the user
	w, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return []Wallet{}, nil
	}

	// 2. Get all active crypto assets from the platform
	assets, _, err := s.platformService.GetCryptoAssets(ctx, true, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("fetch crypto assets: %w", err)
	}

	wallets := make([]Wallet, 0, len(assets))

	// 3. For each asset, fetch balance from blockchain and sync
	for _, asset := range assets {
		// Try to get balance from DB first
		wb, err := s.repo.GetByUserAndCrypto(ctx, userID, asset.Symbol)
		balance := 0.0
		lockedBalance := 0.0
		if err == nil && wb != nil {
			balance = wb.Balance
			lockedBalance = wb.LockedBalance
		}

		// Fetch live balance from blockchain using shared address
		realBalance, err := s.walletClient.GetBalance(ctx, w.Address)
		if err == nil {
			// Update DB if balance changed
			if realBalance != balance {
				balance = realBalance
				// Create virtual wallet to update
				v := Wallet{
					WalletID:      w.WalletID,
					CryptoID:      &asset.ID,
					Balance:       balance,
					LockedBalance: lockedBalance,
				}
				_ = s.repo.Update(ctx, &v)
			}
		}

		// Fetch pending deposits
		pendingBalance, _ := s.repo.GetPendingDepositTotal(ctx, userID, asset.ID)

		// Create a virtual wallet representation for this asset
		virtualWallet := Wallet{
			WalletID:       w.WalletID,
			UserID:         w.UserID,
			CryptoID:       &asset.ID,
			CryptoSymbol:   asset.Symbol,
			Address:        w.Address,
			Balance:        balance,
			LockedBalance:  lockedBalance,
			PendingBalance: pendingBalance,
			IsActive:       w.IsActive,
			IsPrimary:      asset.Symbol == "ETH",
			LastUpdated:    time.Now(),
			CreatedAt:      w.CreatedAt,
		}
		wallets = append(wallets, virtualWallet)
	}

	return wallets, nil
}

func (s *service) CreateDefaultWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error) { // Check if user already has a wallet
	existing, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// 1. Generate a real ECDSA key pair for the wallet
	address, privateKey, err := s.walletClient.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generate key pair: %w", err)
	}

	// 2. Encrypt the private key using the platform secret
	encryptionKey := []byte(s.cfg.WalletEncryptionKey)
	encryptedPk, err := sharedcrypto.Encrypt(privateKey, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt private key: %w", err)
	}

	wallet := &Wallet{
		WalletID:            uuid.New(),
		UserID:              userID,
		Address:             address,
		EncryptedPrivateKey: &encryptedPk,
		IsActive:            true,
		LastUpdated:         time.Now(),
		CreatedAt:           time.Now(),
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		// Check for unique constraint violation - race condition handling
		// If another request created the wallet, return the existing one
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			existing, getErr := s.repo.GetByUser(ctx, userID)
			if getErr == nil && existing != nil {
				return existing, nil
			}
			return nil, apperrors.WalletExists("user already has a wallet", err)
		}
		return nil, fmt.Errorf("create wallet: %w", err)
	}

	return wallet, nil
}

func (s *service) GetDepositAddress(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error) {
	// Single wallet design: everyone has one wallet address that holds all tokens
	wallet, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, apperrors.NotFound("wallet not found for user", nil)
	}

	// Attach the requested symbol to the virtual wallet response
	wallet.CryptoSymbol = strings.TrimSpace(strings.ToUpper(cryptoSymbol))
	return wallet, nil
}

func (s *service) Withdraw(ctx context.Context, userID uuid.UUID, cryptoSymbol, destination string, amount float64, fee float64, memo *string, emailCode string) (*WalletTransaction, error) {
	if amount <= 0 {
		return nil, apperrors.InvalidInput("amount must be greater than zero", nil)
	}
	if strings.TrimSpace(destination) == "" {
		return nil, apperrors.InvalidInput("destination is required", nil)
	}
	if fee < 0 {
		return nil, apperrors.InvalidInput("fee cannot be negative", nil)
	}

	// FR-304: Email confirmation logic (Mocked for now)
	// In production, we would verify emailCode against a stored code in Redis/DB
	if emailCode == "" {
		return nil, apperrors.InvalidInput("email confirmation code is required", nil)
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

func (s *service) GetDailyLimitInfo(ctx context.Context, userID uuid.UUID) (float64, float64, error) {
	dailyTotal, err := s.repo.GetDailyWithdrawalTotal(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	return dailyTotal, 500.0, nil // 500.0 is the daily limit
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
