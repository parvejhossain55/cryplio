package wallet

import (
	"context"
	"database/sql"
	"fmt"

	domainwallet "cryplio/internal/domain/wallet"

	"github.com/google/uuid"
)

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) domainwallet.Repository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetByID(ctx context.Context, walletID uuid.UUID) (*domainwallet.Wallet, error) {
	query := `
		SELECT wallet_id, user_id, crypto_id, address, balance, locked_balance,
		       is_active, false AS is_primary, updated_at AS last_updated, created_at
		FROM wallets
		WHERE wallet_id = $1
	`
	var w domainwallet.Wallet
	if err := r.db.QueryRowContext(ctx, query, walletID).Scan(
		&w.WalletID, &w.UserID, &w.CryptoID, &w.Address, &w.Balance, &w.LockedBalance,
		&w.IsActive, &w.IsPrimary, &w.LastUpdated, &w.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get wallet by id: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) GetByUserAndCrypto(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*domainwallet.Wallet, error) {
	// First try with crypto_assets join for full info
	query := `
		SELECT w.wallet_id, w.user_id, w.crypto_id, w.address, w.balance, w.locked_balance,
		       w.is_active, false AS is_primary, w.updated_at AS last_updated, w.created_at
		FROM wallets w
		JOIN crypto_assets ca ON ca.id = w.crypto_id
		WHERE w.user_id = $1
		  AND UPPER(ca.symbol) = UPPER($2)
		  AND ca.is_active = true
	`
	var w domainwallet.Wallet
	if err := r.db.QueryRowContext(ctx, query, userID, cryptoSymbol).Scan(
		&w.WalletID, &w.UserID, &w.CryptoID, &w.Address, &w.Balance, &w.LockedBalance,
		&w.IsActive, &w.IsPrimary, &w.LastUpdated, &w.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get wallet by user and crypto: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) GetByUserAndCryptoID(ctx context.Context, userID uuid.UUID, cryptoID int) (*domainwallet.Wallet, error) {
	query := `
		SELECT wallet_id, user_id, crypto_id, address, balance, locked_balance,
		       is_active, false AS is_primary, updated_at AS last_updated, created_at
		FROM wallets
		WHERE user_id = $1 AND crypto_id = $2
	`
	var w domainwallet.Wallet
	if err := r.db.QueryRowContext(ctx, query, userID, cryptoID).Scan(
		&w.WalletID, &w.UserID, &w.CryptoID, &w.Address, &w.Balance, &w.LockedBalance,
		&w.IsActive, &w.IsPrimary, &w.LastUpdated, &w.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get wallet by user and crypto_id: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) GetCryptoIDBySymbol(ctx context.Context, symbol string) (int, error) {
	query := `
		SELECT id FROM crypto_assets
		WHERE UPPER(symbol) = UPPER($1) AND is_active = true
		ORDER BY id ASC
		LIMIT 1
	`
	var cryptoID int
	if err := r.db.QueryRowContext(ctx, query, symbol).Scan(&cryptoID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("cryptocurrency not found: %s", symbol)
		}
		return 0, fmt.Errorf("get crypto_id by symbol: %w", err)
	}
	return cryptoID, nil
}

func (r *walletRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domainwallet.Wallet, error) {
	query := `
		SELECT w.wallet_id, w.user_id, w.crypto_id, w.address, w.balance, w.locked_balance,
		       w.is_active, false AS is_primary, w.updated_at AS last_updated, w.created_at,
		       COALESCE(ca.symbol, 'USDT') AS crypto_symbol
		FROM wallets w
		LEFT JOIN crypto_assets ca ON ca.id = w.crypto_id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list wallets by user: %w", err)
	}
	defer rows.Close()

	wallets := make([]domainwallet.Wallet, 0)
	for rows.Next() {
		var w domainwallet.Wallet
		if err := rows.Scan(
			&w.WalletID, &w.UserID, &w.CryptoID, &w.Address, &w.Balance, &w.LockedBalance,
			&w.IsActive, &w.IsPrimary, &w.LastUpdated, &w.CreatedAt, &w.CryptoSymbol,
		); err != nil {
			return nil, fmt.Errorf("scan wallet: %w", err)
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (r *walletRepository) Create(ctx context.Context, wallet *domainwallet.Wallet) error {
	query := `
		INSERT INTO wallets (wallet_id, user_id, crypto_id, address, balance, locked_balance,
		                     is_active, address_label, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	if err := r.db.QueryRowContext(ctx, query,
		wallet.WalletID, wallet.UserID, wallet.CryptoID, wallet.Address,
		wallet.Balance, wallet.LockedBalance, wallet.IsActive, wallet.Address,
	).Scan(&wallet.CreatedAt, &wallet.LastUpdated); err != nil {
		return fmt.Errorf("create wallet: %w", err)
	}
	return nil
}

func (r *walletRepository) Update(ctx context.Context, wallet *domainwallet.Wallet) error {
	query := `
		UPDATE wallets
		SET balance = $1,
		    locked_balance = $2,
		    is_active = $3,
		    updated_at = NOW()
		WHERE wallet_id = $4
		RETURNING updated_at
	`
	if err := r.db.QueryRowContext(ctx, query, wallet.Balance, wallet.LockedBalance, wallet.IsActive, wallet.WalletID).
		Scan(&wallet.LastUpdated); err != nil {
		return fmt.Errorf("update wallet: %w", err)
	}
	return nil
}

func (r *walletRepository) CreateTransaction(ctx context.Context, tx *domainwallet.WalletTransaction) error {
	query := `
		INSERT INTO wallet_transactions (
			txn_id, wallet_id, user_id, type, amount, fee, status, reference_id,
			tx_hash, confirmations, crypto_id, to_address, metadata, created_at
		)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8,
		       $9, $10, w.crypto_id, $11, NULL, NOW()
		FROM wallets w
		WHERE w.wallet_id = $2
		RETURNING created_at
	`
	if err := r.db.QueryRowContext(
		ctx, query,
		tx.TxID, tx.WalletID, tx.UserID, tx.Type, tx.Amount, tx.Fee, tx.Status,
		tx.ReferenceID, tx.TxHash, tx.Confirmations, tx.Memo,
	).Scan(&tx.CreatedAt); err != nil {
		return fmt.Errorf("create wallet transaction: %w", err)
	}
	tx.UpdatedAt = tx.CreatedAt
	return nil
}

func (r *walletRepository) GetTransactionByID(ctx context.Context, txID uuid.UUID) (*domainwallet.WalletTransaction, error) {
	query := `
		SELECT txn_id, wallet_id, user_id, type, status, amount, fee,
		       (amount - fee) AS net_amount, 0::numeric AS balance_after,
		       reference_id, tx_hash, confirmations, to_address, created_at
		FROM wallet_transactions
		WHERE txn_id = $1
	`
	var tx domainwallet.WalletTransaction
	if err := r.db.QueryRowContext(ctx, query, txID).Scan(
		&tx.TxID, &tx.WalletID, &tx.UserID, &tx.Type, &tx.Status, &tx.Amount, &tx.Fee,
		&tx.NetAmount, &tx.BalanceAfter, &tx.ReferenceID, &tx.TxHash, &tx.Confirmations, &tx.Memo, &tx.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get wallet transaction by id: %w", err)
	}
	tx.UpdatedAt = tx.CreatedAt
	return &tx, nil
}

func (r *walletRepository) ListTransactionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domainwallet.WalletTransaction, int, error) {
	countQuery := `SELECT COUNT(*) FROM wallet_transactions WHERE user_id = $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count wallet transactions: %w", err)
	}

	query := `
		SELECT txn_id, wallet_id, user_id, type, status, amount, fee,
		       (amount - fee) AS net_amount, 0::numeric AS balance_after,
		       reference_id, tx_hash, confirmations, to_address, created_at
		FROM wallet_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list wallet transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]domainwallet.WalletTransaction, 0)
	for rows.Next() {
		var tx domainwallet.WalletTransaction
		if err := rows.Scan(
			&tx.TxID, &tx.WalletID, &tx.UserID, &tx.Type, &tx.Status, &tx.Amount, &tx.Fee,
			&tx.NetAmount, &tx.BalanceAfter, &tx.ReferenceID, &tx.TxHash, &tx.Confirmations, &tx.Memo, &tx.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan wallet transaction: %w", err)
		}
		tx.UpdatedAt = tx.CreatedAt
		transactions = append(transactions, tx)
	}

	return transactions, total, nil
}

func (r *walletRepository) GetDailyWithdrawalTotal(ctx context.Context, userID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM wallet_transactions
		WHERE user_id = $1
		  AND type = 'withdrawal'
		  AND status IN ('pending', 'confirmed', 'completed')
		  AND DATE(created_at) = CURRENT_DATE
	`
	var total float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get daily withdrawal total: %w", err)
	}
	return total, nil
}

func (r *walletRepository) ListPendingWithdrawals(ctx context.Context, limit, offset int) ([]domainwallet.WalletTransaction, int, error) {
	countQuery := `SELECT COUNT(*) FROM wallet_transactions WHERE type = 'withdrawal' AND status = 'pending' AND requires_approval = true`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count pending withdrawals: %w", err)
	}

	query := `
		SELECT txn_id, wallet_id, user_id, type, status, amount, fee,
		       (amount - fee) AS net_amount, 0::numeric AS balance_after,
		       reference_id, tx_hash, confirmations, to_address, requires_approval, approved_by, approved_at, created_at
		FROM wallet_transactions
		WHERE type = 'withdrawal' AND status = 'pending' AND requires_approval = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list pending withdrawals: %w", err)
	}
	defer rows.Close()

	transactions := make([]domainwallet.WalletTransaction, 0)
	for rows.Next() {
		var tx domainwallet.WalletTransaction
		if err := rows.Scan(
			&tx.TxID, &tx.WalletID, &tx.UserID, &tx.Type, &tx.Status, &tx.Amount, &tx.Fee,
			&tx.NetAmount, &tx.BalanceAfter, &tx.ReferenceID, &tx.TxHash, &tx.Confirmations, &tx.Memo,
			&tx.RequiresApproval, &tx.ApprovedBy, &tx.ApprovedAt, &tx.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan pending withdrawal: %w", err)
		}
		tx.UpdatedAt = tx.CreatedAt
		transactions = append(transactions, tx)
	}
	return transactions, total, nil
}

func (r *walletRepository) ApproveWithdrawal(ctx context.Context, txID uuid.UUID, adminID uuid.UUID, txHash string) error {
	query := `
		UPDATE wallet_transactions
		SET status = 'completed', tx_hash = $1, approved_by = $2, approved_at = NOW(), updated_at = NOW()
		WHERE txn_id = $3 AND type = 'withdrawal' AND status = 'pending' AND requires_approval = true
	`
	res, err := r.db.ExecContext(ctx, query, txHash, adminID, txID)
	if err != nil {
		return fmt.Errorf("approve withdrawal: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("approve withdrawal rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("withdrawal not found or already processed")
	}
	return nil
}

func (r *walletRepository) RejectWithdrawal(ctx context.Context, txID uuid.UUID, adminID uuid.UUID, reason string) error {
	query := `
		UPDATE wallet_transactions
		SET status = 'cancelled', approved_by = $1, approved_at = NOW(), updated_at = NOW(),
		    metadata = COALESCE(metadata, '{}') || jsonb_build_object('rejection_reason', $2)
		WHERE txn_id = $3 AND type = 'withdrawal' AND status = 'pending' AND requires_approval = true
	`
	res, err := r.db.ExecContext(ctx, query, adminID, reason, txID)
	if err != nil {
		return fmt.Errorf("reject withdrawal: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("reject withdrawal rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("withdrawal not found or already processed")
	}
	return nil
}
