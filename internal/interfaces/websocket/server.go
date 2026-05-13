package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	sharedjwt "cryplio/pkg/jwt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    string      `json:"user_id,omitempty"`
	TradeID   string      `json:"trade_id,omitempty"`
}

// Client represents a WebSocket client
type Client struct {
	ID       string
	UserID   uuid.UUID
	Conn     *websocket.Conn
	Send     chan Message
	Server   *Server
	TradeIDs map[string]bool // Trades this client is subscribed to
}

// Server represents the WebSocket server
type Server struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	mutex      sync.RWMutex
	upgrader   websocket.Upgrader
	jwtSecret  string
}

// NewServer creates a new WebSocket server
func NewServer(jwtSecret string) *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message, 256),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		jwtSecret: jwtSecret,
	}
}

// Start starts the WebSocket server
func (s *Server) Start(ctx context.Context) {
	go s.run()

	http.HandleFunc("/ws", s.handleWebSocket)

	log.Println("WebSocket server started on /ws")
}

// run handles the main server loop
func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
			log.Printf("Client connected: %s (User: %s)", client.ID, client.UserID)

			// Send welcome message
			select {
			case client.Send <- Message{
				Type:      "connected",
				Data:      map[string]string{"message": "Connected to Cryplio WebSocket"},
				Timestamp: time.Now(),
			}:
			default:
				close(client.Send)
				delete(s.clients, client)
			}

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.Send)
				log.Printf("Client disconnected: %s (User: %s)", client.ID, client.UserID)
			}
			s.mutex.Unlock()

		case message := <-s.broadcast:
			s.mutex.RLock()
			for client := range s.clients {
				// Filter messages based on subscriptions
				if s.shouldSendToClient(client, message) {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(s.clients, client)
					}
				}
			}
			s.mutex.RUnlock()
		}
	}
}

// shouldSendToClient determines if a message should be sent to a client
func (s *Server) shouldSendToClient(client *Client, message Message) bool {
	// Send if no specific trade ID or if client is subscribed to the trade
	if message.TradeID == "" {
		return true
	}

	_, subscribed := client.TradeIDs[message.TradeID]
	return subscribed
}

// HandleWebSocket exposes the WebSocket handler for external routers (public)
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	s.handleWebSocket(w, r)
}

// handleWebSocket handles new WebSocket connections (internal)
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Extract user ID from query parameter or JWT token
	userIDStr := r.URL.Query().Get("user_id")
	tokenStr := r.URL.Query().Get("token")

	var userID uuid.UUID

	if userIDStr != "" {
		var parseErr error
		userID, parseErr = uuid.Parse(userIDStr)
		if parseErr != nil {
			log.Printf("Invalid user ID parameter: %v", parseErr)
			conn.Close()
			return
		}
	} else if tokenStr != "" {
		claims, parseErr := sharedjwt.Parse(s.jwtSecret, tokenStr)
		if parseErr != nil {
			log.Printf("Invalid WebSocket token: %v", parseErr)
			conn.Close()
			return
		}

		uid, ok := claims[sharedjwt.ClaimUserID].(string)
		if !ok {
			log.Printf("User ID not found in token")
			conn.Close()
			return
		}

		var uuidErr error
		userID, uuidErr = uuid.Parse(uid)
		if uuidErr != nil {
			log.Printf("Invalid user ID in token: %v", uuidErr)
			conn.Close()
			return
		}
	} else {
		log.Printf("WebSocket connection attempted without user_id or token")
		conn.Close()
		return
	}

	client := &Client{
		ID:       uuid.New().String(),
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan Message, 256),
		Server:   s,
		TradeIDs: make(map[string]bool),
	}

	client.Server.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump handles messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.Server.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	for {
		var msg json.RawMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message and handle subscriptions
		var message struct {
			Type    string          `json:"type"`
			Data    json.RawMessage `json:"data"`
			TradeID string          `json:"trade_id,omitempty"`
		}

		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Message parse error: %v", err)
			continue
		}

		// Handle subscription messages
		switch message.Type {
		case "subscribe_trade":
			if message.TradeID != "" {
				c.TradeIDs[message.TradeID] = true
				log.Printf("Client %s subscribed to trade %s", c.ID, message.TradeID)
			}
		case "unsubscribe_trade":
			if message.TradeID != "" {
				delete(c.TradeIDs, message.TradeID)
				log.Printf("Client %s unsubscribed from trade %s", c.ID, message.TradeID)
			}
		}
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("Write error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// BroadcastMessage broadcasts a message to all connected clients
func (s *Server) BroadcastMessage(messageType string, data interface{}, tradeID string) {
	message := Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
		TradeID:   tradeID,
	}

	select {
	case s.broadcast <- message:
	default:
		log.Printf("Broadcast channel full, dropping message")
	}
}

// BroadcastToUser sends a message to a specific user
func (s *Server) BroadcastToUser(userID uuid.UUID, messageType string, data interface{}) {
	message := Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID.String(),
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for client := range s.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(s.clients, client)
			}
		}
	}
}

// GetConnectedUsers returns the list of connected user IDs
func (s *Server) GetConnectedUsers() []uuid.UUID {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]uuid.UUID, 0, len(s.clients))
	seen := make(map[uuid.UUID]bool)

	for client := range s.clients {
		if !seen[client.UserID] {
			users = append(users, client.UserID)
			seen[client.UserID] = true
		}
	}

	return users
}
