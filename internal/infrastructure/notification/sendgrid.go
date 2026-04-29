package notification

import "context"

// EmailMessage describes an outbound email notification.
type EmailMessage struct {
	To      string
	Subject string
	HTML    string
	Text    string
}

// SendGridClient sends transactional email notifications.
type SendGridClient interface {
	Send(ctx context.Context, message EmailMessage) error
}

// NoopSendGridClient is a placeholder until the provider integration is built.
type NoopSendGridClient struct{}

func NewSendGridClient() *NoopSendGridClient {
	return &NoopSendGridClient{}
}

func (c *NoopSendGridClient) Send(context.Context, EmailMessage) error {
	return nil
}
