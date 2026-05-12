package wallet

import "context"

// WalletClient handles deposit address provisioning, sends, and chain watchers.
type WalletClient interface {
	GenerateKeyPair() (address, privateKey string, err error)
	CreateDepositAddress(ctx context.Context, cryptoID int, userID string) (string, error)
	GetBalance(ctx context.Context, address string) (float64, error)
	Send(ctx context.Context, tx *WalletTransaction, destination string) (string, error)
	Watch(ctx context.Context, txHash string) error
}
