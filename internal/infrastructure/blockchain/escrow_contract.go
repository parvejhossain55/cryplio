package blockchain

import (
	"context"

	trading "cryplio/internal/domain/trading"
)

// EscrowContractClient coordinates escrow lifecycle actions against a chain adapter.
type EscrowContractClient interface {
	Lock(ctx context.Context, trade *trading.Trade) (txHash string, contractAddress string, err error)
	Release(ctx context.Context, trade *trading.Trade) (txHash string, err error)
	Refund(ctx context.Context, trade *trading.Trade) (txHash string, err error)
}

// NoopEscrowContractClient is a stub until an EVM client is integrated.
type NoopEscrowContractClient struct{}

func NewEscrowContractClient() *NoopEscrowContractClient {
	return &NoopEscrowContractClient{}
}

func (c *NoopEscrowContractClient) Lock(context.Context, *trading.Trade) (string, string, error) {
	return "", "", nil
}

func (c *NoopEscrowContractClient) Release(context.Context, *trading.Trade) (string, error) {
	return "", nil
}

func (c *NoopEscrowContractClient) Refund(context.Context, *trading.Trade) (string, error) {
	return "", nil
}
