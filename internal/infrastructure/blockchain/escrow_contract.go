package blockchain

import (
	"context"
	trading "cryplio/internal/domain/trading"
)

// NoopEscrowContractClient is a stub until an EVM client is integrated.
type NoopEscrowContractClient struct{}

func NewNoopEscrowContractClient() *NoopEscrowContractClient {
	return &NoopEscrowContractClient{}
}

func (c *NoopEscrowContractClient) Lock(context.Context, *trading.Trade) (string, string, error) {
	return "0xnoop-lock", "0xnoop-contract", nil
}

func (c *NoopEscrowContractClient) Release(context.Context, *trading.Trade) (string, error) {
	return "0xnoop-release", nil
}

func (c *NoopEscrowContractClient) Refund(context.Context, *trading.Trade) (string, error) {
	return "0xnoop-refund", nil
}

func (c *NoopEscrowContractClient) AdminRefund(context.Context, *trading.Trade) (string, error) {
	return "0xnoop-admin-refund", nil
}
