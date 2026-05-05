package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, walletID uuid.UUID) (*Wallet, error)
	GetByUserAndCrypto(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Wallet, error)
	Update(ctx context.Context, wallet *Wallet) error
	CreateTransaction(ctx context.Context, tx *WalletTransaction) error
	GetTransactionByID(ctx context.Context, txID uuid.UUID) (*WalletTransaction, error)
	ListTransactionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]WalletTransaction, int, error)
}
