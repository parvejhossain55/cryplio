package wallet

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"cryplio/internal/infrastructure/blockchain"
	"cryplio/pkg/database"

	"github.com/google/uuid"
)

// WalletService handles wallet operations
type WalletService struct {
	db         *database.DB
	usdtClient *blockchain.USDTWalletClient
	escrowAddr string
}

// NewWalletService creates a new wallet service
func NewWalletService(db *database.DB, usdtClient *blockchain.USDTWalletClient, escrowAddr string) *WalletService {
	return &WalletService{
		db:         db,
		usdtClient: usdtClient,
		escrowAddr: escrowAddr,
	}
}

// Wallet represents a user's cryptocurrency wallet
type Wallet struct {
	WalletID       string  `json:"wallet_id"`
	UserID         string  `json:"user_id"`
	CryptoSymbol   string  `json:"crypto_symbol"`
	BlockchainAddr string  `json:"blockchain_address"`
	Balance        float64 `json:"balance"`
	LockedBalance  float64 `json:"locked_balance"`
	IsActive       bool    `json:"is_active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	TxID             string  `json:"tx_id"`
	WalletID         string  `json:"wallet_id"`
	Type             string  `json:"type"`
	Amount           float64 `json:"amount"`
	BalanceBefore    float64 `json:"balance_before"`
	BalanceAfter     float64 `json:"balance_after"`
	Status           string  `json:"status"`
	BlockchainTxHash string  `json:"blockchain_tx_hash"`
	ToAddress        string  `json:"to_address"`
	FromAddress      string  `json:"from_address"`
	NetworkFee       float64 `json:"network_fee"`
	TradeID          string  `json:"trade_id"`
	TwoFAVerified    bool    `json:"two_fa_verified"`
	Notes            string  `json:"notes"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// WithdrawalRequest represents a withdrawal request
type WithdrawalRequest struct {
	RequestID        string  `json:"request_id"`
	UserID           string  `json:"user_id"`
	WalletID         string  `json:"wallet_id"`
	Amount           float64 `json:"amount"`
	ToAddress        string  `json:"to_address"`
	BlockchainTxHash string  `json:"blockchain_tx_hash"`
	Status           string  `json:"status"`
	ApprovedBy       string  `json:"approved_by"`
	ApprovedAt       string  `json:"approved_at"`
	RejectedBy       string  `json:"rejected_by"`
	RejectedAt       string  `json:"rejected_at"`
	RejectionReason  string  `json:"rejection_reason"`
	TwoFACode        string  `json:"two_fa_code"`
	Notes            string  `json:"notes"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// CreateWallet creates a new wallet for a user
func (s *WalletService) CreateWallet(ctx context.Context, userID string, cryptoSymbol string) (*Wallet, error) {
	walletID := uuid.New().String()
	blockchainAddr := s.usdtClient.GetAddress()
	now := time.Now().UTC().Format(time.RFC3339)

	query := `
		INSERT INTO wallets (wallet_id, user_id, crypto_symbol, blockchain_address, balance, locked_balance, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING wallet_id, user_id, crypto_symbol, blockchain_address, balance, locked_balance, is_active, created_at, updated_at
	`

	var wallet Wallet
	err := s.db.QueryRowContext(ctx, query,
		walletID, userID, cryptoSymbol, blockchainAddr, 0.0, 0.0, true, now, now,
	).Scan(
		&wallet.WalletID, &wallet.UserID, &wallet.CryptoSymbol, &wallet.BlockchainAddr,
		&wallet.Balance, &wallet.LockedBalance, &wallet.IsActive, &wallet.CreatedAt, &wallet.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create wallet: %w", err)
	}

	return &wallet, nil
}

// GetWalletByUser retrieves a wallet by user ID and crypto symbol
func (s *WalletService) GetWalletByUser(ctx context.Context, userID string, cryptoSymbol string) (*Wallet, error) {
	query := `
		SELECT wallet_id, user_id, crypto_symbol, blockchain_address, balance, locked_balance, is_active, created_at, updated_at
		FROM wallets
		WHERE user_id = $1 AND crypto_symbol = $2
	`

	var wallet Wallet
	err := s.db.QueryRowContext(ctx, query, userID, cryptoSymbol).Scan(
		&wallet.WalletID, &wallet.UserID, &wallet.CryptoSymbol, &wallet.BlockchainAddr,
		&wallet.Balance, &wallet.LockedBalance, &wallet.IsActive, &wallet.CreatedAt, &wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("get wallet: %w", err)
	}

	return &wallet, nil
}

// GetUserWallets retrieves all wallets for a user
func (s *WalletService) GetUserWallets(ctx context.Context, userID string) ([]Wallet, error) {
	query := `
		SELECT wallet_id, user_id, crypto_symbol, blockchain_address, balance, locked_balance, is_active, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query user wallets: %w", err)
	}
	defer rows.Close()

	var wallets []Wallet
	for rows.Next() {
		var wallet Wallet
		err := rows.Scan(
			&wallet.WalletID, &wallet.UserID, &wallet.CryptoSymbol, &wallet.BlockchainAddr,
			&wallet.Balance, &wallet.LockedBalance, &wallet.IsActive, &wallet.CreatedAt, &wallet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wallet: %w", err)
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// UpdateWalletBalance updates a wallet's balance
func (s *WalletService) UpdateWalletBalance(ctx context.Context, walletID string, newBalance float64, newLockedBalance float64) error {
	now := time.Now().UTC().Format(time.RFC3339)

	query := `
		UPDATE wallets
		SET balance = $2, locked_balance = $3, updated_at = $4
		WHERE wallet_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, walletID, newBalance, newLockedBalance, now)
	if err != nil {
		return fmt.Errorf("update wallet balance: %w", err)
	}

	return nil
}

// CreateWalletTransaction creates a new wallet transaction
func (s *WalletService) CreateWalletTransaction(ctx context.Context, tx *WalletTransaction) error {
	now := time.Now().UTC().Format(time.RFC3339)

	query := `
		INSERT INTO wallet_transactions (
			tx_id, wallet_id, type, amount, balance_before, balance_after, status,
			blockchain_tx_hash, to_address, from_address, network_fee, trade_id,
			two_fa_verified, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := s.db.ExecContext(ctx, query,
		tx.TxID, tx.WalletID, tx.Type, tx.Amount, tx.BalanceBefore, tx.BalanceAfter, tx.Status,
		tx.BlockchainTxHash, tx.ToAddress, tx.FromAddress, tx.NetworkFee, tx.TradeID,
		tx.TwoFAVerified, tx.Notes, now, now,
	)

	if err != nil {
		return fmt.Errorf("create wallet transaction: %w", err)
	}

	return nil
}

// GetWalletTransactions retrieves transactions for a wallet
func (s *WalletService) GetWalletTransactions(ctx context.Context, walletID string, limit int) ([]WalletTransaction, error) {
	query := `
		SELECT tx_id, wallet_id, type, amount, balance_before, balance_after, status,
			   blockchain_tx_hash, to_address, from_address, network_fee, trade_id,
			   two_fa_verified, notes, created_at, updated_at
		FROM wallet_transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.QueryContext(ctx, query, walletID, limit)
	if err != nil {
		return nil, fmt.Errorf("query wallet transactions: %w", err)
	}
	defer rows.Close()

	var transactions []WalletTransaction
	for rows.Next() {
		var tx WalletTransaction
		err := rows.Scan(
			&tx.TxID, &tx.WalletID, &tx.Type, &tx.Amount, &tx.BalanceBefore, &tx.BalanceAfter, &tx.Status,
			&tx.BlockchainTxHash, &tx.ToAddress, &tx.FromAddress, &tx.NetworkFee, &tx.TradeID,
			&tx.TwoFAVerified, &tx.Notes, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wallet transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// ProcessDeposit processes a USDT deposit
func (s *WalletService) ProcessDeposit(ctx context.Context, userID string, amount float64, txHash string) error {
	// Get or create USDT wallet for user
	wallet, err := s.GetWalletByUser(ctx, userID, "USDT")
	if err != nil {
		if err.Error() == "wallet not found" {
			wallet, err = s.CreateWallet(ctx, userID, "USDT")
			if err != nil {
				return fmt.Errorf("create wallet for deposit: %w", err)
			}
		} else {
			return fmt.Errorf("get wallet for deposit: %w", err)
		}
	}

	// Verify transaction on blockchain
	successful, err := s.usdtClient.GetTransactionStatus(ctx, txHash)
	if err != nil {
		return fmt.Errorf("verify transaction: %w", err)
	}

	if !successful {
		return fmt.Errorf("transaction failed or not found")
	}

	// Update wallet balance
	newBalance := wallet.Balance + amount
	err = s.UpdateWalletBalance(ctx, wallet.WalletID, newBalance, wallet.LockedBalance)
	if err != nil {
		return fmt.Errorf("update wallet balance for deposit: %w", err)
	}

	// Create transaction record
	tx := &WalletTransaction{
		TxID:             uuid.New().String(),
		WalletID:         wallet.WalletID,
		Type:             "deposit",
		Amount:           amount,
		BalanceBefore:    wallet.Balance,
		BalanceAfter:     newBalance,
		Status:           "confirmed",
		BlockchainTxHash: txHash,
		ToAddress:        wallet.BlockchainAddr,
		TwoFAVerified:    false,
		Notes:            "USDT deposit",
	}

	err = s.CreateWalletTransaction(ctx, tx)
	if err != nil {
		log.Printf("Warning: failed to create deposit transaction record: %v", err)
	}

	return nil
}

// CreateWithdrawalRequest creates a withdrawal request
func (s *WalletService) CreateWithdrawalRequest(ctx context.Context, userID string, amount float64, toAddress string, twoFACode string) (*WithdrawalRequest, error) {
	// Get user's USDT wallet
	wallet, err := s.GetWalletByUser(ctx, userID, "USDT")
	if err != nil {
		return nil, fmt.Errorf("get wallet for withdrawal: %w", err)
	}

	// Check if user has sufficient balance
	if wallet.Balance < amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create withdrawal request
	requestID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	query := `
		INSERT INTO withdrawal_requests (
			request_id, user_id, wallet_id, amount, to_address, status, two_fa_code, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING request_id, user_id, wallet_id, amount, to_address, blockchain_tx_hash, status,
				  approved_by, approved_at, rejected_by, rejected_at, rejection_reason, two_fa_code, notes, created_at, updated_at
	`

	var request WithdrawalRequest
	err = s.db.QueryRowContext(ctx, query,
		requestID, userID, wallet.WalletID, amount, toAddress, "pending", twoFACode, "User withdrawal request", now, now,
	).Scan(
		&request.RequestID, &request.UserID, &request.WalletID, &request.Amount, &request.ToAddress,
		&request.BlockchainTxHash, &request.Status, &request.ApprovedBy, &request.ApprovedAt,
		&request.RejectedBy, &request.RejectedAt, &request.RejectionReason,
		&request.TwoFACode, &request.Notes, &request.CreatedAt, &request.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create withdrawal request: %w", err)
	}

	return &request, nil
}

// ProcessWithdrawal processes an approved withdrawal
func (s *WalletService) ProcessWithdrawal(ctx context.Context, requestID string) error {
	// Get withdrawal request
	request, err := s.getWithdrawalRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get withdrawal request: %w", err)
	}

	if request.Status != "approved" {
		return fmt.Errorf("withdrawal request not approved")
	}

	// Get wallet
	wallet, err := s.GetWalletByUser(ctx, request.UserID, "USDT")
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}

	// Check balance again
	if wallet.Balance < request.Amount {
		return fmt.Errorf("insufficient balance")
	}

	// Send USDT on blockchain
	amountWei := blockchain.ToUSDTAmount(request.Amount)
	txHash, err := s.usdtClient.Transfer(ctx, request.ToAddress, amountWei)
	if err != nil {
		return fmt.Errorf("send USDT: %w", err)
	}

	// Lock funds temporarily
	newBalance := wallet.Balance - request.Amount
	newLockedBalance := wallet.LockedBalance + request.Amount
	err = s.UpdateWalletBalance(ctx, wallet.WalletID, newBalance, newLockedBalance)
	if err != nil {
		return fmt.Errorf("lock funds: %w", err)
	}

	// Update withdrawal request
	err = s.updateWithdrawalRequest(ctx, requestID, "processing", txHash, "", "", "", "")
	if err != nil {
		return fmt.Errorf("update withdrawal request: %w", err)
	}

	// Create transaction record
	tx := &WalletTransaction{
		TxID:             uuid.New().String(),
		WalletID:         wallet.WalletID,
		Type:             "withdrawal",
		Amount:           request.Amount,
		BalanceBefore:    wallet.Balance,
		BalanceAfter:     newBalance,
		Status:           "pending",
		BlockchainTxHash: txHash,
		ToAddress:        request.ToAddress,
		FromAddress:      wallet.BlockchainAddr,
		TwoFAVerified:    true,
		Notes:            "USDT withdrawal",
	}

	err = s.CreateWalletTransaction(ctx, tx)
	if err != nil {
		log.Printf("Warning: failed to create withdrawal transaction record: %v", err)
	}

	// Wait for transaction confirmation in background
	go s.confirmWithdrawal(context.Background(), requestID, txHash)

	return nil
}

// confirmWithdrawal confirms a withdrawal transaction
func (s *WalletService) confirmWithdrawal(ctx context.Context, requestID string, txHash string) {
	// Wait for transaction to be mined
	receipt, err := s.usdtClient.WaitForTransaction(ctx, txHash)
	if err != nil {
		log.Printf("Failed to wait for withdrawal transaction %s: %v", txHash, err)
		s.updateWithdrawalRequest(ctx, requestID, "failed", txHash, "", "", "Transaction failed", "")
		return
	}

	if receipt.Status == 1 {
		// Transaction successful - unlock funds
		err = s.finalizeWithdrawal(ctx, requestID)
		if err != nil {
			log.Printf("Failed to finalize withdrawal %s: %v", requestID, err)
		}
	} else {
		// Transaction failed - refund locked funds
		err = s.refundWithdrawal(ctx, requestID)
		if err != nil {
			log.Printf("Failed to refund withdrawal %s: %v", requestID, err)
		}
	}
}

// finalizeWithdrawal finalizes a successful withdrawal
func (s *WalletService) finalizeWithdrawal(ctx context.Context, requestID string) error {
	request, err := s.getWithdrawalRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get withdrawal request: %w", err)
	}

	wallet, err := s.GetWalletByUser(ctx, request.UserID, "USDT")
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}

	// Remove locked funds
	newLockedBalance := wallet.LockedBalance - request.Amount
	err = s.UpdateWalletBalance(ctx, wallet.WalletID, wallet.Balance, newLockedBalance)
	if err != nil {
		return fmt.Errorf("update wallet balance: %w", err)
	}

	// Update withdrawal request
	return s.updateWithdrawalRequest(ctx, requestID, "completed", request.BlockchainTxHash, "", "", "", "")
}

// refundWithdrawal refunds a failed withdrawal
func (s *WalletService) refundWithdrawal(ctx context.Context, requestID string) error {
	request, err := s.getWithdrawalRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get withdrawal request: %w", err)
	}

	wallet, err := s.GetWalletByUser(ctx, request.UserID, "USDT")
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}

	// Refund locked funds
	newBalance := wallet.Balance + request.Amount
	newLockedBalance := wallet.LockedBalance - request.Amount
	err = s.UpdateWalletBalance(ctx, wallet.WalletID, newBalance, newLockedBalance)
	if err != nil {
		return fmt.Errorf("update wallet balance: %w", err)
	}

	// Update withdrawal request
	return s.updateWithdrawalRequest(ctx, requestID, "failed", request.BlockchainTxHash, "", "", "Transaction failed", "")
}

// Helper functions
func (s *WalletService) getWithdrawalRequest(ctx context.Context, requestID string) (*WithdrawalRequest, error) {
	query := `
		SELECT request_id, user_id, wallet_id, amount, to_address, blockchain_tx_hash, status,
			   approved_by, approved_at, rejected_by, rejected_at, rejection_reason, two_fa_code, notes, created_at, updated_at
		FROM withdrawal_requests
		WHERE request_id = $1
	`

	var request WithdrawalRequest
	err := s.db.QueryRowContext(ctx, query, requestID).Scan(
		&request.RequestID, &request.UserID, &request.WalletID, &request.Amount, &request.ToAddress,
		&request.BlockchainTxHash, &request.Status, &request.ApprovedBy, &request.ApprovedAt,
		&request.RejectedBy, &request.RejectedAt, &request.RejectionReason,
		&request.TwoFACode, &request.Notes, &request.CreatedAt, &request.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get withdrawal request: %w", err)
	}

	return &request, nil
}

func (s *WalletService) updateWithdrawalRequest(ctx context.Context, requestID string, status string, txHash string, approvedBy string, rejectedBy string, rejectionReason string, notes string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	query := `
		UPDATE withdrawal_requests
		SET status = $2, blockchain_tx_hash = $3, approved_by = $4, rejected_by = $5,
			rejection_reason = $6, notes = $7, updated_at = $8
		WHERE request_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, requestID, status, txHash, approvedBy, rejectedBy, rejectionReason, notes, now)
	if err != nil {
		return fmt.Errorf("update withdrawal request: %w", err)
	}

	return nil
}
