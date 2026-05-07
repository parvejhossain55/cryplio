package payment

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// MockPaymentProvider provides mock implementations for various payment providers
type MockPaymentProvider struct {
	name         string
	transactions map[string]*MockTransaction
	mutex        sync.RWMutex
}

// MockTransaction represents a mock payment transaction
type MockTransaction struct {
	ID            string            `json:"id"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Status        string            `json:"status"` // "pending", "completed", "failed", "cancelled"
	Provider      string            `json:"provider"`
	Type          string            `json:"type"` // "deposit", "withdrawal", "payment"
	AccountID     string            `json:"account_id"`
	Reference     string            `json:"reference"`
	Description   string            `json:"description"`
	Metadata      map[string]string `json:"metadata"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty"`
	FailureReason string            `json:"failure_reason,omitempty"`
}

// PaymentProvider interface
type PaymentProvider interface {
	Name() string
	ProcessPayment(ctx context.Context, req *MockPaymentRequest) (*PaymentResponse, error)
	ProcessWithdrawal(ctx context.Context, req *MockWithdrawalRequest) (*PaymentResponse, error)
	GetTransactionStatus(ctx context.Context, transactionID string) (*TransactionStatus, error)
	GetBalance(ctx context.Context, accountID string) (*Balance, error)
}

// MockPaymentRequest extends PaymentRequest with additional fields
type MockPaymentRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	AccountID   string            `json:"account_id"`
	Reference   string            `json:"reference"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// MockWithdrawalRequest represents a withdrawal request
type MockWithdrawalRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	AccountID   string            `json:"account_id"`
	Reference   string            `json:"reference"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	ProcessedAt   time.Time `json:"processed_at"`
}

// TransactionStatus represents transaction status
type TransactionStatus struct {
	ID            string     `json:"id"`
	Status        string     `json:"status"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	FailureReason string     `json:"failure_reason,omitempty"`
}

// Balance represents account balance
type Balance struct {
	AccountID string             `json:"account_id"`
	Available map[string]float64 `json:"available"`
	Reserved  map[string]float64 `json:"reserved"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// NewMockBkashProvider creates a mock Bkash provider
func NewMockBkashProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		name:         "bkash",
		transactions: make(map[string]*MockTransaction),
	}
}

// NewMockNagadProvider creates a mock Nagad provider
func NewMockNagadProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		name:         "nagad",
		transactions: make(map[string]*MockTransaction),
	}
}

// NewMockBankTransferProvider creates a mock bank transfer provider
func NewMockBankTransferProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		name:         "bank_transfer",
		transactions: make(map[string]*MockTransaction),
	}
}

// NewMockWiseProvider creates a mock Wise provider
func NewMockWiseProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		name:         "wise",
		transactions: make(map[string]*MockTransaction),
	}
}

// NewMockPayPalProvider creates a mock PayPal provider
func NewMockPayPalProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		name:         "paypal",
		transactions: make(map[string]*MockTransaction),
	}
}

// NewMockUPIProvider creates a mock UPI provider
func NewMockUPIProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		name:         "upi",
		transactions: make(map[string]*MockTransaction),
	}
}

// Name returns the provider name
func (m *MockPaymentProvider) Name() string {
	return m.name
}

// ProcessPayment processes a mock payment
func (m *MockPaymentProvider) ProcessPayment(ctx context.Context, req *MockPaymentRequest) (*PaymentResponse, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	transactionID := fmt.Sprintf("%s_tx_%d", m.name, time.Now().UnixNano())

	transaction := &MockTransaction{
		ID:          transactionID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      "pending",
		Provider:    m.name,
		Type:        "payment",
		AccountID:   req.AccountID,
		Reference:   req.Reference,
		Description: req.Description,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.transactions[transactionID] = transaction

	// Simulate async processing
	go m.simulatePaymentProcessing(transactionID)

	log.Printf("Mock payment initiated: %s for %s, amount: %.2f %s",
		transactionID, m.name, req.Amount, req.Currency)

	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        "pending",
		Message:       "Payment initiated successfully",
		ProcessedAt:   time.Now(),
	}, nil
}

// ProcessWithdrawal processes a mock withdrawal
func (m *MockPaymentProvider) ProcessWithdrawal(ctx context.Context, req *MockWithdrawalRequest) (*PaymentResponse, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	transactionID := fmt.Sprintf("%s_wd_%d", m.name, time.Now().UnixNano())

	transaction := &MockTransaction{
		ID:          transactionID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      "pending",
		Provider:    m.name,
		Type:        "withdrawal",
		AccountID:   req.AccountID,
		Reference:   req.Reference,
		Description: req.Description,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.transactions[transactionID] = transaction

	// Simulate async processing
	go m.simulateWithdrawalProcessing(transactionID)

	log.Printf("Mock withdrawal initiated: %s for %s, amount: %.2f %s",
		transactionID, m.name, req.Amount, req.Currency)

	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        "pending",
		Message:       "Withdrawal initiated successfully",
		ProcessedAt:   time.Now(),
	}, nil
}

// GetTransactionStatus returns transaction status
func (m *MockPaymentProvider) GetTransactionStatus(ctx context.Context, transactionID string) (*TransactionStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	transaction, exists := m.transactions[transactionID]
	if !exists {
		return nil, fmt.Errorf("transaction not found: %s", transactionID)
	}

	return &TransactionStatus{
		ID:            transaction.ID,
		Status:        transaction.Status,
		Amount:        transaction.Amount,
		Currency:      transaction.Currency,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		CompletedAt:   transaction.CompletedAt,
		FailureReason: transaction.FailureReason,
	}, nil
}

// GetBalance returns account balance
func (m *MockPaymentProvider) GetBalance(ctx context.Context, accountID string) (*Balance, error) {
	// Return mock balance
	return &Balance{
		AccountID: accountID,
		Available: map[string]float64{
			"USD": 10000.0 + rand.Float64()*50000,
			"EUR": 8000.0 + rand.Float64()*40000,
			"BDT": 500000.0 + rand.Float64()*2000000,
		},
		Reserved: map[string]float64{
			"USD": rand.Float64() * 1000,
			"EUR": rand.Float64() * 800,
			"BDT": rand.Float64() * 50000,
		},
		UpdatedAt: time.Now(),
	}, nil
}

// simulatePaymentProcessing simulates payment processing
func (m *MockPaymentProvider) simulatePaymentProcessing(transactionID string) {
	// Simulate processing delay
	delay := time.Duration(5+rand.Intn(15)) * time.Second
	time.Sleep(delay)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	transaction, exists := m.transactions[transactionID]
	if !exists {
		return
	}

	// 90% success rate for payments
	if rand.Float32() < 0.9 {
		transaction.Status = "completed"
		completedAt := time.Now()
		transaction.CompletedAt = &completedAt
		log.Printf("Mock payment completed: %s", transactionID)
	} else {
		transaction.Status = "failed"
		transaction.FailureReason = "Insufficient funds"
		log.Printf("Mock payment failed: %s - %s", transactionID, transaction.FailureReason)
	}

	transaction.UpdatedAt = time.Now()
}

// simulateWithdrawalProcessing simulates withdrawal processing
func (m *MockPaymentProvider) simulateWithdrawalProcessing(transactionID string) {
	// Simulate processing delay
	delay := time.Duration(10+rand.Intn(30)) * time.Second
	time.Sleep(delay)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	transaction, exists := m.transactions[transactionID]
	if !exists {
		return
	}

	// 85% success rate for withdrawals
	if rand.Float32() < 0.85 {
		transaction.Status = "completed"
		completedAt := time.Now()
		transaction.CompletedAt = &completedAt
		log.Printf("Mock withdrawal completed: %s", transactionID)
	} else {
		transaction.Status = "failed"
		reasons := []string{"Account verification required", "Daily limit exceeded", "Bank account not found"}
		transaction.FailureReason = reasons[rand.Intn(len(reasons))]
		log.Printf("Mock withdrawal failed: %s - %s", transactionID, transaction.FailureReason)
	}

	transaction.UpdatedAt = time.Now()
}

// GetAllTransactions returns all transactions for testing
func (m *MockPaymentProvider) GetAllTransactions() []*MockTransaction {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	transactions := make([]*MockTransaction, 0, len(m.transactions))
	for _, tx := range m.transactions {
		txCopy := *tx
		transactions = append(transactions, &txCopy)
	}

	return transactions
}

// GetStats returns provider statistics
func (m *MockPaymentProvider) GetStats() map[string]int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]int{
		"total":     0,
		"pending":   0,
		"completed": 0,
		"failed":    0,
		"cancelled": 0,
	}

	for _, tx := range m.transactions {
		stats["total"]++
		stats[tx.Status]++
	}

	return stats
}
