package blockchain

import (
	"context"
	wallet "cryplio/internal/domain/wallet"
)

// NoopWalletClient is a stub until a chain wallet provider is integrated.
type NoopWalletClient struct{}

func NewNoopWalletClient() *NoopWalletClient {
	return &NoopWalletClient{}
}

func (c *NoopWalletClient) CreateDepositAddress(context.Context, int, string) (string, error) {
	return "0xnoop-address", nil
}

func (c *NoopWalletClient) Send(context.Context, *wallet.WalletTransaction, string) (string, error) {
	return "0xnoop-tx", nil
}

func (c *NoopWalletClient) Watch(context.Context, string) error {
	return nil
}
