package ws

import "context"

// Message is the websocket payload envelope for chat and notification streams.
type Message struct {
	Channel string
	UserID  string
	Event   string
	Data    []byte
}

// Hub coordinates websocket fan-out for live delivery concerns.
type Hub interface {
	Publish(ctx context.Context, message Message) error
	Run(ctx context.Context) error
}

// NoopHub is a compile-safe placeholder until websocket delivery is wired in.
type NoopHub struct{}

func NewHub() *NoopHub {
	return &NoopHub{}
}

func (h *NoopHub) Publish(context.Context, Message) error {
	return nil
}

func (h *NoopHub) Run(context.Context) error {
	return nil
}
