package payment

import "context"

// PaymentRequest captures the minimum fields shared by fiat payment adapters.
type PaymentRequest struct {
	Reference string
	Amount    float64
	Currency  string
	Account   string
}

// PaymentResult captures the provider reference for a payment action.
type PaymentResult struct {
	ProviderReference string
	Status            string
}

// BKashClient handles bKash-specific payment operations.
type BKashClient interface {
	CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error)
}

// NoopBKashClient is a placeholder until the provider integration is built.
type NoopBKashClient struct{}

func NewBKashClient() *NoopBKashClient {
	return &NoopBKashClient{}
}

func (c *NoopBKashClient) CreatePayment(context.Context, PaymentRequest) (*PaymentResult, error) {
	return nil, nil
}
