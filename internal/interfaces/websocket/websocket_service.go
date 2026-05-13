package websocket

import (
	"context"
	"net/http"
	"sync"

	"cryplio/pkg/logger"

	"github.com/google/uuid"
)

type websocketService struct {
	server *Server
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewService creates a new WebSocket service
func NewService(jwtSecret string) Service {
	return &websocketService{
		server: NewServer(jwtSecret),
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

	logger.Info("WebSocket service started", logger.Fields{})
	return nil
}

// Stop stops the WebSocket service
func (ws *websocketService) Stop() error {
	if ws.cancel != nil {
		ws.cancel()
	}
	ws.wg.Wait()
	logger.Info("WebSocket service stopped", logger.Fields{})
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

// HandleWebSocket handles WebSocket upgrade requests
func (ws *websocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws.server.HandleWebSocket(w, r)
}
