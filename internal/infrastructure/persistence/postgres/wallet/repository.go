package wallet

// repository.go is the entry point for the wallet Postgres repository.
// Method implementations are split across focused files:
//
//   wallet.go      — Wallet CRUD (balances, addresses, crypto lookups)
//   transaction.go — WalletTransaction CRUD (deposits, withdrawals, approvals)

import (
	"database/sql"

	domainwallet "cryplio/internal/domain/wallet"
)

// walletRepository implements domainwallet.Repository on top of PostgreSQL.
type walletRepository struct {
	db *sql.DB
}

// NewWalletRepository constructs a walletRepository backed by the given *sql.DB.
func NewWalletRepository(db *sql.DB) domainwallet.Repository {
	return &walletRepository{db: db}
}
