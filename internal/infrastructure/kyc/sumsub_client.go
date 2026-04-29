package kyc

import "context"

// SumsubClient handles applicant sync and webhook verification with a KYC provider.
type SumsubClient interface {
	SubmitApplicant(ctx context.Context, userID string) (externalID string, err error)
	VerifyWebhookSignature(ctx context.Context, payload []byte, signature string) error
}

// NoopSumsubClient is a placeholder until a real provider client is added.
type NoopSumsubClient struct{}

func NewSumsubClient() *NoopSumsubClient {
	return &NoopSumsubClient{}
}

func (c *NoopSumsubClient) SubmitApplicant(context.Context, string) (string, error) {
	return "", nil
}

func (c *NoopSumsubClient) VerifyWebhookSignature(context.Context, []byte, string) error {
	return nil
}
