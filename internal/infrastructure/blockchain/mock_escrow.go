package blockchain

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cryplio/internal/domain/trading"

	"github.com/google/uuid"
)

// MockEscrowContractClient provides a mock implementation of the escrow contract
type MockEscrowContractClient struct {
	escrows map[string]*MockEscrow
	mutex   sync.RWMutex
}

// MockEscrow represents a mock escrow record
type MockEscrow struct {
	ID           string    `json:"id"`
	TradeID      uuid.UUID `json:"trade_id"`
	Buyer        string    `json:"buyer"`
	Seller       string    `json:"seller"`
	Amount       string    `json:"amount"`
	Token        string    `json:"token"`
	Status       string    `json:"status"` // "created", "locked", "released", "refunded"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at";` // Added semicolon here
	ReleasedTo   string    `json:"released_to,omitempty"`
	RefundReason string    `json:"refund_reason,omitempty"`
}

// NewMockEscrowContractClient creates a new mock escrow client
func NewMockEscrowContractClient() *MockEscrowContractClient {
	return &MockEscrowContractClient{
		escrows: make(map[string]*MockEscrow),
	}
}

// CreateEscrow creates a new mock escrow
func (m *MockEscrowContractClient) CreateEscrow(ctx context.Context, tradeID uuid.UUID, buyer, seller, amount, token string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	escrowID := fmt.Sprintf("mock_escrow_%s_%d", tradeID.String()[:8], time.Now().Unix())

	escrow := &MockEscrow{
		ID:        escrowID,
		TradeID:   tradeID,
		Buyer:     buyer,
		Seller:    seller,
		Amount:    amount,
		Token:     token,
		Status:    "created",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.escrows[escrowID] = escrow

	log.Printf("Mock escrow created: %s for trade %s", escrowID, tradeID)

	return escrowID, nil
}

// Lock locks funds in the mock escrow
func (m *MockEscrowContractClient) Lock(ctx context.Context, trade *trading.Trade) (txHash string, contractAddress string, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	escrowID := fmt.Sprintf("mock_escrow_%s", trade.TradeID.String())

	// Create escrow if it doesn't exist
	if _, exists := m.escrows[escrowID]; !exists {
		m.escrows[escrowID] = &MockEscrow{
			ID:        escrowID,
			TradeID:   trade.TradeID,
			Buyer:     trade.BuyerID.String(),
			Seller:    trade.SellerID.String(),
			Amount:    fmt.Sprintf("%.6f", trade.CryptoAmount),
			Token:     "USDT", // Default token for MVP
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	escrow := m.escrows[escrowID]
	if escrow.Status != "created" {
		return "", "", fmt.Errorf("escrow cannot be locked, current status: %s", escrow.Status)
	}

	escrow.Status = "locked"
	escrow.UpdatedAt = time.Now()

	txHash = fmt.Sprintf("0xmocktx_%d", time.Now().UnixNano())
	contractAddress = "0xmock_contract_address"

	log.Printf("Mock escrow locked: %s, tx: %s", escrowID, txHash)

	return txHash, contractAddress, nil
}

// Release releases funds to the buyer in the mock escrow
func (m *MockEscrowContractClient) Release(ctx context.Context, trade *trading.Trade) (txHash string, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	escrowID := fmt.Sprintf("mock_escrow_%s", trade.TradeID.String())
	escrow, exists := m.escrows[escrowID]
	if !exists {
		return "", fmt.Errorf("escrow not found: %s", escrowID)
	}

	if escrow.Status != "locked" {
		return "", fmt.Errorf("escrow cannot be released, current status: %s", escrow.Status)
	}

	escrow.Status = "released"
	escrow.ReleasedTo = trade.BuyerID.String()
	escrow.UpdatedAt = time.Now()

	txHash = fmt.Sprintf("0xmocktx_%d", time.Now().UnixNano())

	log.Printf("Mock escrow released: %s, tx: %s", escrowID, txHash)

	return txHash, nil
}

// Refund refunds funds to the seller in the mock escrow
func (m *MockEscrowContractClient) Refund(ctx context.Context, trade *trading.Trade) (txHash string, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	escrowID := fmt.Sprintf("mock_escrow_%s", trade.TradeID.String())
	escrow, exists := m.escrows[escrowID]
	if !exists {
		return "", fmt.Errorf("escrow not found: %s", escrowID)
	}

	if escrow.Status != "locked" {
		return "", fmt.Errorf("escrow cannot be refunded, current status: %s", escrow.Status)
	}

	escrow.Status = "refunded"
	escrow.RefundReason = "Trade dispute or cancellation"
	escrow.UpdatedAt = time.Now()

	txHash = fmt.Sprintf("0xmocktx_%d", time.Now().UnixNano())

	log.Printf("Mock escrow refunded: %s, tx: %s", escrowID, txHash)

	return txHash, nil
}

// GetEscrowStatus returns the status of a mock escrow
func (m *MockEscrowContractClient) GetEscrowStatus(ctx context.Context, escrowID string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	escrow, exists := m.escrows[escrowID]
	if !exists {
		return "", fmt.Errorf("escrow not found: %s", escrowID)
	}

	return escrow.Status, nil
}

// GetEscrowDetails returns the details of a mock escrow
func (m *MockEscrowContractClient) GetEscrowDetails(ctx context.Context, escrowID string) (*MockEscrow, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	escrow, exists := m.escrows[escrowID]
	if !exists {
		return nil, fmt.Errorf("escrow not found: %s", escrowID)
	}

	// Return a copy to avoid concurrent modification
	escrowCopy := *escrow
	return &escrowCopy, nil
}

// ListEscrows returns all mock escrows
func (m *MockEscrowContractClient) ListEscrows(ctx context.Context) ([]*MockEscrow, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	escrows := make([]*MockEscrow, 0, len(m.escrows))
	for _, escrow := range m.escrows {
		escrowCopy := *escrow
		escrows = append(escrows, &escrowCopy)
	}

	return escrows, nil
}

// GetEscrowsByTradeID returns escrows for a specific trade
func (m *MockEscrowContractClient) GetEscrowsByTradeID(ctx context.Context, tradeID uuid.UUID) ([]*MockEscrow, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var escrows []*MockEscrow
	for _, escrow := range m.escrows {
		if escrow.TradeID == tradeID {
			escrowCopy := *escrow
			escrows = append(escrows, &escrowCopy)
		}
	}

	return escrows, nil
}

// SimulateBlockchainDelay simulates blockchain transaction delay
func (m *MockEscrowContractClient) SimulateBlockchainDelay() {
	time.Sleep(time.Duration(100+time.Now().UnixNano()%400) * time.Millisecond)
}

// CleanupExpiredEscrows removes old mock escrows (for testing)
func (m *MockEscrowContractClient) CleanupExpiredEscrows(ctx context.Context, olderThan time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	cutoff := time.Now().Add(-olderThan)
	for id, escrow := range m.escrows {
		if escrow.CreatedAt.Before(cutoff) &&
			(escrow.Status == "released" || escrow.Status == "refunded") {
			delete(m.escrows, id)
			log.Printf("Cleaned up expired mock escrow: %s", id)
		}
	}
}

// GetStats returns mock escrow statistics
func (m *MockEscrowContractClient) GetStats(ctx context.Context) map[string]int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]int{
		"total":    0,
		"created":  0,
		"locked":   0,
		"released": 0,
		"refunded": 0,
	}

	for _, escrow := range m.escrows {
		stats["total"]++
		stats[escrow.Status]++
	}

	return stats
}
