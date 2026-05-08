package websocket

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Service defines the WebSocket service interface
type Service interface {
	Start(ctx context.Context) error
	Stop() error
	BroadcastMessage(messageType string, data interface{}, tradeID string)
	BroadcastToUser(userID uuid.UUID, messageType string, data interface{})
	GetConnectedUsers() []uuid.UUID
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
}

// NotificationEvent represents different types of notifications
type NotificationEvent struct {
	Type      string      `json:"type"`
	UserID    uuid.UUID   `json:"user_id,omitempty"`
	TradeID   uuid.UUID   `json:"trade_id,omitempty"`
	Title     string      `json:"title"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        string `json:"id"`
	TradeID   string `json:"trade_id"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	FileURL   string `json:"file_url,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
	FileSize  int    `json:"file_size,omitempty"`
	CreatedAt string `json:"created_at"`
}

// TradeUpdate represents a trade status update
type TradeUpdate struct {
	TradeID   uuid.UUID `json:"trade_id"`
	Status    string    `json:"status"`
	BuyerID   uuid.UUID `json:"buyer_id"`
	SellerID  uuid.UUID `json:"seller_id"`
	Amount    float64   `json:"amount"`
	Crypto    string    `json:"crypto"`
	Fiat      string    `json:"fiat"`
	Timestamp string    `json:"timestamp"`
}

// MarketUpdate represents market price updates
type MarketUpdate struct {
	CryptoSymbol string  `json:"crypto_symbol"`
	FiatSymbol   string  `json:"fiat_symbol"`
	Price        float64 `json:"price"`
	Change24h    float64 `json:"change_24h"`
	Timestamp    string  `json:"timestamp"`
}
