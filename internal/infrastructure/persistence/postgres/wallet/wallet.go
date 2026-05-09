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
	wallet_id, user_id, crypto_id, address, balance, locked_balance,
	is_active, false AS is_primary, updated_at AS last_updated, created_at
`

// scanWallet scans the standard wallet projection into w.
func scanWallet(row interface{ Scan(...any) error }, w *domainwallet.Wallet) error {
	var cryptoID sql.NullInt64
	err := row.Scan(
		&w.WalletID, &w.UserID, &cryptoID, &w.Address, &w.Balance, &w.LockedBalance,
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
		SELECT w.wallet_id, w.user_id, w.crypto_id, w.address, w.balance, w.locked_balance,
		       w.is_active, false AS is_primary, w.updated_at AS last_updated, w.created_at
		FROM wallets w
		JOIN crypto_assets ca ON ca.id = w.crypto_id
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
		SELECT w.wallet_id, w.user_id, w.crypto_id, w.address, w.balance, w.locked_balance,
		       w.is_active, false AS is_primary, w.updated_at AS last_updated, w.created_at,
		       ca.symbol AS crypto_symbol
		FROM wallets w
		LEFT JOIN crypto_assets ca ON ca.id = w.crypto_id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC`, userID)
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
			&w.WalletID, &w.UserID, &cryptoID, &w.Address, &w.Balance, &w.LockedBalance,
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
			wallet_id, user_id, crypto_id, address, balance, locked_balance,
			is_active, address_label, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING created_at, updated_at`,
		wallet.WalletID, wallet.UserID, cryptoID, wallet.Address,
		wallet.Balance, wallet.LockedBalance, wallet.IsActive, wallet.Address,
	).Scan(&wallet.CreatedAt, &wallet.LastUpdated)
}

func (r *walletRepository) Update(ctx context.Context, wallet *domainwallet.Wallet) error {
	return r.db.QueryRowContext(ctx, `
		UPDATE wallets
		SET balance = $1, locked_balance = $2, is_active = $3, updated_at = NOW()
		WHERE wallet_id = $4
		RETURNING updated_at`,
		wallet.Balance, wallet.LockedBalance, wallet.IsActive, wallet.WalletID,
	).Scan(&wallet.LastUpdated)
}
