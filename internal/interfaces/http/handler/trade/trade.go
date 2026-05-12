package trade

// trade.go defines the TradeHandler and its constructor.
// Handler implementations are split across focused files:
//
//   ad.go        — Trade advertisement handlers
//   lifecycle.go — Trade status/lifecycle handlers
//   chat.go      — Chat message and feedback handlers

import (
	"cryplio/internal/domain/trading"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/interfaces/websocket"
)

// TradeHandler handles all trade-related HTTP endpoints.
type TradeHandler struct {
	adManager     trading.AdManager
	lifecycle     trading.TradeLifecycleManager
	communication trading.TradeCommunicationManager
	storage       storage.ObjectStorage
	wsService     websocket.Service
}

// NewTradeHandler creates a new TradeHandler.
func NewTradeHandler(
	adManager trading.AdManager,
	lifecycle trading.TradeLifecycleManager,
	communication trading.TradeCommunicationManager,
	storage storage.ObjectStorage,
	wsService websocket.Service,
) *TradeHandler {
	return &TradeHandler{
		adManager:     adManager,
		lifecycle:     lifecycle,
		communication: communication,
		storage:       storage,
		wsService:     wsService,
	}
}
