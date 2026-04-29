package payment

import "context"

// WiseClient handles Wise payout and transfer operations.
type WiseClient interface {
	CreateTransfer(ctx context.Context, req PaymentRequest) (*PaymentResult, error)
}

// NoopWiseClient is a placeholder until the provider integration is built.
type NoopWiseClient struct{}

func NewWiseClient() *NoopWiseClient {
	return &NoopWiseClient{}
}

func (c *NoopWiseClient) CreateTransfer(context.Context, PaymentRequest) (*PaymentResult, error) {
	return nil, nil
}
