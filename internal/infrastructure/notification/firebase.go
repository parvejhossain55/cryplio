package notification

import "context"

// PushMessage describes an outbound mobile push notification.
type PushMessage struct {
	DeviceToken string
	Title       string
	Body        string
}

// FirebaseClient sends push notifications.
type FirebaseClient interface {
	SendPush(ctx context.Context, message PushMessage) error
}

// NoopFirebaseClient is a placeholder until the provider integration is built.
type NoopFirebaseClient struct{}

func NewFirebaseClient() *NoopFirebaseClient {
	return &NoopFirebaseClient{}
}

func (c *NoopFirebaseClient) SendPush(context.Context, PushMessage) error {
	return nil
}
