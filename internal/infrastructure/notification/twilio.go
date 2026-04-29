package notification

import "context"

// SMSMessage describes an outbound SMS notification.
type SMSMessage struct {
	To   string
	Body string
}

// TwilioClient sends SMS notifications.
type TwilioClient interface {
	SendSMS(ctx context.Context, message SMSMessage) error
}

// NoopTwilioClient is a placeholder until the provider integration is built.
type NoopTwilioClient struct{}

func NewTwilioClient() *NoopTwilioClient {
	return &NoopTwilioClient{}
}

func (c *NoopTwilioClient) SendSMS(context.Context, SMSMessage) error {
	return nil
}
