package websocket

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

type websocketService struct {
	server *Server
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewService creates a new WebSocket service
func NewService() Service {
	return &websocketService{
		server: NewServer(),
	}
}

// Start starts the WebSocket service
func (ws *websocketService) Start(ctx context.Context) error {
	ws.ctx, ws.cancel = context.WithCancel(ctx)
	
	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()
		ws.server.Start(ws.ctx)
	}()
	
	log.Println("WebSocket service started")
	return nil
}

// Stop stops the WebSocket service
func (ws *websocketService) Stop() error {
	if ws.cancel != nil {
		ws.cancel()
	}
	ws.wg.Wait()
	log.Println("WebSocket service stopped")
	return nil
}

// BroadcastMessage broadcasts a message to all connected clients
func (ws *websocketService) BroadcastMessage(messageType string, data interface{}, tradeID string) {
	ws.server.BroadcastMessage(messageType, data, tradeID)
}

// BroadcastToUser sends a message to a specific user
func (ws *websocketService) BroadcastToUser(userID uuid.UUID, messageType string, data interface{}) {
	ws.server.BroadcastToUser(userID, messageType, data)
}

// GetConnectedUsers returns the list of connected user IDs
func (ws *websocketService) GetConnectedUsers() []uuid.UUID {
	return ws.server.GetConnectedUsers()
}
