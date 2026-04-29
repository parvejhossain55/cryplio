package events

import "context"

// Event represents a domain event.
type Event interface {
	// Name returns the event name/type.
	Name() string
}

// Handler handles an event.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Dispatcher dispatches events to registered handlers.
type Dispatcher interface {
	Dispatch(ctx context.Context, event Event) error
	Register(name string, handler Handler) error
}

// NoopDispatcher is a no-operation dispatcher (default for no events).
type NoopDispatcher struct{}

func NewNoopDispatcher() *NoopDispatcher {
	return &NoopDispatcher{}
}

func (d *NoopDispatcher) Dispatch(ctx context.Context, event Event) error {
	// no-op
	return nil
}

func (d *NoopDispatcher) Register(name string, handler Handler) error {
	// no-op
	return nil
}
