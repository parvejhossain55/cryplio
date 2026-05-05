package wallet

import "context"

// WalletClient handles deposit address provisioning, sends, and chain watchers.
type WalletClient interface {
	CreateDepositAddress(ctx context.Context, cryptoID int, userID string) (string, error)
	Send(ctx context.Context, tx *WalletTransaction, destination string) (string, error)
	Watch(ctx context.Context, txHash string) error
}
