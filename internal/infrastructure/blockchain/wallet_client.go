package blockchain

import (
	"context"

	wallet "cryplio/internal/domain/wallet"
)

// WalletClient handles deposit address provisioning, sends, and chain watchers.
type WalletClient interface {
	CreateDepositAddress(ctx context.Context, cryptoID int, userID string) (string, error)
	Send(ctx context.Context, tx *wallet.WalletTransaction, destination string) (string, error)
	Watch(ctx context.Context, txHash string) error
}

// NoopWalletClient is a stub until a chain wallet provider is integrated.
type NoopWalletClient struct{}

func NewWalletClient() *NoopWalletClient {
	return &NoopWalletClient{}
}

func (c *NoopWalletClient) CreateDepositAddress(context.Context, int, string) (string, error) {
	return "", nil
}

func (c *NoopWalletClient) Send(context.Context, *wallet.WalletTransaction, string) (string, error) {
	return "", nil
}

func (c *NoopWalletClient) Watch(context.Context, string) error {
	return nil
}
