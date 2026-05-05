package trading

import "context"

// EscrowContractClient coordinates escrow lifecycle actions against a chain adapter.
type EscrowContractClient interface {
	Lock(ctx context.Context, trade *Trade) (txHash string, contractAddress string, err error)
	Release(ctx context.Context, trade *Trade) (txHash string, err error)
	Refund(ctx context.Context, trade *Trade) (txHash string, err error)
}
