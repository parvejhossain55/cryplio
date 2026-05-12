package wallet

import (
	"context"
	"database/sql"
	"fmt"

	domainwallet "cryplio/internal/domain/wallet"

	"github.com/google/uuid"
)

// ─── Wallet CRUD ──────────────────────────────────────────────────────────────

// walletColumns is the canonical SELECT for a wallet row.
const walletColumns = `
	wallet_id, user_id, crypto_id, address, encrypted_private_key, balance, locked_balance,
	is_active, false AS is_primary, updated_at AS last_updated, created_at
`

// scanWallet scans the standard wallet projection into w.
func scanWallet(row interface{ Scan(...any) error }, w *domainwallet.Wallet) error {
	var cryptoID sql.NullInt64
	err := row.Scan(
		&w.WalletID, &w.UserID, &cryptoID, &w.Address, &w.EncryptedPrivateKey, &w.Balance, &w.LockedBalance,
		&w.IsActive, &w.IsPrimary, &w.LastUpdated, &w.CreatedAt,
	)
	if err != nil {
		return err
	}
	if cryptoID.Valid {
		id := int(cryptoID.Int64)
		w.CryptoID = &id
	}
	return nil
}

func (r *walletRepository) GetByID(ctx context.Context, walletID uuid.UUID) (*domainwallet.Wallet, error) {
	var w domainwallet.Wallet
	err := scanWallet(r.db.QueryRowContext(ctx,
		`SELECT `+walletColumns+` FROM wallets WHERE wallet_id = $1`, walletID), &w)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get wallet by id: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) GetByUser(ctx context.Context, userID uuid.UUID) (*domainwallet.Wallet, error) {
	var w domainwallet.Wallet
	err := scanWallet(r.db.QueryRowContext(ctx,
		`SELECT `+walletColumns+` FROM wallets WHERE user_id = $1`, userID), &w)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get wallet by user: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) GetByUserAndCrypto(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*domainwallet.Wallet, error) {
	var w domainwallet.Wallet
	err := scanWallet(r.db.QueryRowContext(ctx, `
		SELECT w.wallet_id, w.user_id, wb.crypto_id, w.address, w.encrypted_private_key, wb.balance, wb.locked_balance,
		       w.is_active, false AS is_primary, wb.updated_at AS last_updated, w.created_at
		FROM wallets w
		JOIN wallet_balances wb ON wb.wallet_id = w.wallet_id
		JOIN crypto_assets ca ON ca.id = wb.crypto_id
		WHERE w.user_id = $1 AND UPPER(ca.symbol) = UPPER($2) AND ca.is_active = true`,
		userID, cryptoSymbol), &w)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get wallet by user and crypto: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) GetByUserAndCryptoID(ctx context.Context, userID uuid.UUID, cryptoID int) (*domainwallet.Wallet, error) {
	var w domainwallet.Wallet
	err := scanWallet(r.db.QueryRowContext(ctx,
		`SELECT `+walletColumns+` FROM wallets WHERE user_id = $1 AND crypto_id = $2`,
		userID, cryptoID), &w)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get wallet by user and crypto_id: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) GetCryptoIDBySymbol(ctx context.Context, symbol string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM crypto_assets
		WHERE UPPER(symbol) = UPPER($1) AND is_active = true
		ORDER BY id ASC LIMIT 1`, symbol).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("cryptocurrency not found: %s", symbol)
	}
	if err != nil {
		return 0, fmt.Errorf("get crypto_id by symbol: %w", err)
	}
	return id, nil
}

func (r *walletRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domainwallet.Wallet, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT w.wallet_id, w.user_id, wb.crypto_id, w.address, w.encrypted_private_key, wb.balance, wb.locked_balance,
		       w.is_active, false AS is_primary, wb.updated_at AS last_updated, w.created_at,
		       ca.symbol AS crypto_symbol
		FROM wallets w
		JOIN wallet_balances wb ON wb.wallet_id = w.wallet_id
		JOIN crypto_assets ca ON ca.id = wb.crypto_id
		WHERE w.user_id = $1
		ORDER BY ca.symbol ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list wallets by user: %w", err)
	}
	defer rows.Close()

	wallets := make([]domainwallet.Wallet, 0)
	for rows.Next() {
		var w domainwallet.Wallet
		var cryptoID sql.NullInt64
		var cryptoSymbol sql.NullString
		if err := rows.Scan(
			&w.WalletID, &w.UserID, &cryptoID, &w.Address, &w.EncryptedPrivateKey, &w.Balance, &w.LockedBalance,
			&w.IsActive, &w.IsPrimary, &w.LastUpdated, &w.CreatedAt, &cryptoSymbol,
		); err != nil {
			return nil, fmt.Errorf("scan wallet: %w", err)
		}
		if cryptoID.Valid {
			id := int(cryptoID.Int64)
			w.CryptoID = &id
		}
		if cryptoSymbol.Valid {
			w.CryptoSymbol = cryptoSymbol.String
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (r *walletRepository) Create(ctx context.Context, wallet *domainwallet.Wallet) error {
	var cryptoID sql.NullInt64
	if wallet.CryptoID != nil {
		cryptoID.Int64 = int64(*wallet.CryptoID)
		cryptoID.Valid = true
	}
	return r.db.QueryRowContext(ctx, `
		INSERT INTO wallets (
			wallet_id, user_id, crypto_id, address, encrypted_private_key, balance, locked_balance,
			is_active, address_label, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at`,
		wallet.WalletID, wallet.UserID, cryptoID, wallet.Address, wallet.EncryptedPrivateKey,
		wallet.Balance, wallet.LockedBalance, wallet.IsActive, wallet.Address,
	).Scan(&wallet.CreatedAt, &wallet.LastUpdated)
}

func (r *walletRepository) Update(ctx context.Context, wallet *domainwallet.Wallet) error {
	if wallet.CryptoID == nil {
		return fmt.Errorf("crypto_id is required to update balance")
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO wallet_balances (wallet_id, crypto_id, balance, locked_balance, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (wallet_id, crypto_id) DO UPDATE 
		SET balance = EXCLUDED.balance, 
		    locked_balance = EXCLUDED.locked_balance,
		    updated_at = NOW()`,
		wallet.WalletID, *wallet.CryptoID, wallet.Balance, wallet.LockedBalance)
	return err
}

func (r *walletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, balance float64) error {
	// Note: This method is now ambiguous without cryptoID.
	// In a single-wallet multi-token system, we should ideally use UpdateBalanceByAsset.
	// For compatibility, we'll assume the primary asset (e.g., crypto_id=1) or handle it elsewhere.
	// However, it's better to update it to include cryptoID.
	return fmt.Errorf("UpdateBalance is deprecated, use Update instead with cryptoID")
}
