package payment

import "context"

// NagadClient handles Nagad-specific payment operations.
type NagadClient interface {
	CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error)
}

// NoopNagadClient is a placeholder until the provider integration is built.
type NoopNagadClient struct{}

func NewNagadClient() *NoopNagadClient {
	return &NoopNagadClient{}
}

func (c *NoopNagadClient) CreatePayment(context.Context, PaymentRequest) (*PaymentResult, error) {
	return nil, nil
}
