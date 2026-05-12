package wallet

import (
	"context"
	"database/sql"
	"fmt"

	domainwallet "cryplio/internal/domain/wallet"

	"github.com/google/uuid"
)

// ─── Wallet Transactions ──────────────────────────────────────────────────────

// txColumns is the standard SELECT projection for wallet_transactions.
const txColumns = `
	t.txn_id, t.wallet_id, t.user_id, t.type, t.status, t.amount, t.fee,
	(t.amount - t.fee) AS net_amount, 0::numeric AS balance_after,
	t.reference_id, t.tx_hash, t.confirmations, t.to_address, t.created_at,
	ca.symbol AS crypto_symbol
`

func scanTx(row interface{ Scan(...any) error }, tx *domainwallet.WalletTransaction) error {
	err := row.Scan(
		&tx.TxID, &tx.WalletID, &tx.UserID, &tx.Type, &tx.Status, &tx.Amount, &tx.Fee,
		&tx.NetAmount, &tx.BalanceAfter,
		&tx.ReferenceID, &tx.TxHash, &tx.Confirmations, &tx.Memo, &tx.CreatedAt,
		&tx.CryptoSymbol,
	)
	if err == nil {
		tx.UpdatedAt = tx.CreatedAt
	}
	return err
}

func (r *walletRepository) CreateTransaction(ctx context.Context, tx *domainwallet.WalletTransaction) error {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO wallet_transactions (
			txn_id, wallet_id, user_id, type, amount, fee, status, reference_id,
			tx_hash, confirmations, crypto_id, to_address, metadata, created_at
		)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8,
		       $9, $10, w.crypto_id, $11, NULL, NOW()
		FROM wallets w WHERE w.wallet_id = $2
		RETURNING created_at`,
		tx.TxID, tx.WalletID, tx.UserID, tx.Type, tx.Amount, tx.Fee, tx.Status,
		tx.ReferenceID, tx.TxHash, tx.Confirmations, tx.Memo,
	).Scan(&tx.CreatedAt)
	if err != nil {
		return fmt.Errorf("create wallet transaction: %w", err)
	}
	tx.UpdatedAt = tx.CreatedAt
	return nil
}

func (r *walletRepository) GetTransactionByID(ctx context.Context, txID uuid.UUID) (*domainwallet.WalletTransaction, error) {
	var tx domainwallet.WalletTransaction
	err := scanTx(r.db.QueryRowContext(ctx,
		`SELECT `+txColumns+` 
		 FROM wallet_transactions t
		 JOIN crypto_assets ca ON ca.id = t.crypto_id
		 WHERE t.txn_id = $1`, txID), &tx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get wallet transaction: %w", err)
	}
	return &tx, nil
}

func (r *walletRepository) ListTransactionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domainwallet.WalletTransaction, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM wallet_transactions WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count wallet transactions: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT `+txColumns+`
		 FROM wallet_transactions t
		 JOIN crypto_assets ca ON ca.id = t.crypto_id
		 WHERE t.user_id = $1
		 ORDER BY t.created_at DESC
		 LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list wallet transactions: %w", err)
	}
	defer rows.Close()

	txs := make([]domainwallet.WalletTransaction, 0)
	for rows.Next() {
		var tx domainwallet.WalletTransaction
		if err := scanTx(rows, &tx); err != nil {
			return nil, 0, fmt.Errorf("scan wallet transaction: %w", err)
		}
		txs = append(txs, tx)
	}
	return txs, total, nil
}

func (r *walletRepository) GetPendingDepositTotal(ctx context.Context, userID uuid.UUID, cryptoID int) (float64, error) {
	var total float64
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM wallet_transactions
		WHERE user_id = $1 AND crypto_id = $2 AND type = 'deposit' AND status = 'pending'`,
		userID, cryptoID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get pending deposit total: %w", err)
	}
	return total, nil
}

func (r *walletRepository) GetDailyWithdrawalTotal(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM wallet_transactions
		WHERE user_id = $1
		  AND type = 'withdrawal'
		  AND status IN ('pending', 'confirmed', 'completed')
		  AND DATE(created_at) = CURRENT_DATE`, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get daily withdrawal total: %w", err)
	}
	return total, nil
}

func (r *walletRepository) ListPendingWithdrawals(ctx context.Context, limit, offset int) ([]domainwallet.WalletTransaction, int, error) {
	const pendingWhere = `type = 'withdrawal' AND status = 'pending' AND requires_approval = true`

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM wallet_transactions WHERE `+pendingWhere).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count pending withdrawals: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT txn_id, wallet_id, user_id, type, status, amount, fee,
		       (amount - fee) AS net_amount, 0::numeric AS balance_after,
		       reference_id, tx_hash, confirmations, to_address,
		       requires_approval, approved_by, approved_at, created_at
		FROM wallet_transactions
		WHERE `+pendingWhere+`
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list pending withdrawals: %w", err)
	}
	defer rows.Close()

	txs := make([]domainwallet.WalletTransaction, 0)
	for rows.Next() {
		var tx domainwallet.WalletTransaction
		if err := rows.Scan(
			&tx.TxID, &tx.WalletID, &tx.UserID, &tx.Type, &tx.Status, &tx.Amount, &tx.Fee,
			&tx.NetAmount, &tx.BalanceAfter,
			&tx.ReferenceID, &tx.TxHash, &tx.Confirmations, &tx.Memo,
			&tx.RequiresApproval, &tx.ApprovedBy, &tx.ApprovedAt, &tx.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan pending withdrawal: %w", err)
		}
		tx.UpdatedAt = tx.CreatedAt
		txs = append(txs, tx)
	}
	return txs, total, nil
}

func (r *walletRepository) ApproveWithdrawal(ctx context.Context, txID, adminID uuid.UUID, txHash string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE wallet_transactions
		SET status = 'completed', tx_hash = $1, approved_by = $2,
		    approved_at = NOW(), updated_at = NOW()
		WHERE txn_id = $3 AND type = 'withdrawal' AND status = 'pending' AND requires_approval = true`,
		txHash, adminID, txID)
	if err != nil {
		return fmt.Errorf("approve withdrawal: %w", err)
	}
	return requireOneRow(res, "withdrawal")
}

func (r *walletRepository) RejectWithdrawal(ctx context.Context, txID, adminID uuid.UUID, reason string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE wallet_transactions
		SET status = 'cancelled', approved_by = $1, approved_at = NOW(), updated_at = NOW(),
		    metadata = COALESCE(metadata, '{}') || jsonb_build_object('rejection_reason', $2)
		WHERE txn_id = $3 AND type = 'withdrawal' AND status = 'pending' AND requires_approval = true`,
		adminID, reason, txID)
	if err != nil {
		return fmt.Errorf("reject withdrawal: %w", err)
	}
	return requireOneRow(res, "withdrawal")
}

// requireOneRow returns an error when a mutating query affected zero rows.
func requireOneRow(res interface{ RowsAffected() (int64, error) }, entity string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("%s not found or already processed", entity)
	}
	return nil
}
