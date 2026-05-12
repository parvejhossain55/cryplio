package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, walletID uuid.UUID) (*Wallet, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	GetByUserAndCrypto(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*Wallet, error)
	GetByUserAndCryptoID(ctx context.Context, userID uuid.UUID, cryptoID int) (*Wallet, error)
	GetCryptoIDBySymbol(ctx context.Context, symbol string) (int, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Wallet, error)
	Create(ctx context.Context, wallet *Wallet) error
	Update(ctx context.Context, wallet *Wallet) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, balance float64) error
	CreateTransaction(ctx context.Context, tx *WalletTransaction) error
	GetTransactionByID(ctx context.Context, txID uuid.UUID) (*WalletTransaction, error)
	ListTransactionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]WalletTransaction, int, error)
	GetPendingDepositTotal(ctx context.Context, userID uuid.UUID, cryptoID int) (float64, error)
	GetDailyWithdrawalTotal(ctx context.Context, userID uuid.UUID) (float64, error)
	ListPendingWithdrawals(ctx context.Context, limit, offset int) ([]WalletTransaction, int, error)
	ApproveWithdrawal(ctx context.Context, txID, adminID uuid.UUID, txHash string) error
	RejectWithdrawal(ctx context.Context, txID, adminID uuid.UUID, reason string) error
}
