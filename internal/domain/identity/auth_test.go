package identity

import (
	"context"
	"testing"
	"time"

	"cryplio/internal/domain/wallet"

	"github.com/google/uuid"
)

// MockUserRepository implements identity.UserRepository for testing.
type MockUserRepository struct {
	UserRepository // Embed to satisfy interface with panics for unimplemented methods
	users          map[string]*User
	sessions       map[string]*UserSession
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:    make(map[string]*User),
		sessions: make(map[string]*UserSession),
	}
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	return m.users[email], nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) CreateSession(ctx context.Context, session *UserSession) error {
	m.sessions[session.TokenID] = session
	return nil
}

// MockWalletService implements wallet.Service for testing.
type MockWalletService struct {
	// embed if necessary, but we only need CreateDefaultWallet
	CreateCalled bool
	CreatedFor   uuid.UUID
}

func (m *MockWalletService) GetBalances(ctx context.Context, userID uuid.UUID) ([]wallet.Wallet, error) {
	return nil, nil
}
func (m *MockWalletService) GetDepositAddress(ctx context.Context, userID uuid.UUID, cryptoSymbol string) (*wallet.Wallet, error) {
	return nil, nil
}
func (m *MockWalletService) CreateDefaultWallet(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error) {
	m.CreateCalled = true
	m.CreatedFor = userID
	return &wallet.Wallet{WalletID: uuid.New(), UserID: userID}, nil
}
func (m *MockWalletService) Withdraw(ctx context.Context, userID uuid.UUID, cryptoSymbol, destination string, amount float64, fee float64, memo *string, emailCode string) (*wallet.WalletTransaction, error) {
	return nil, nil
}
func (m *MockWalletService) GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]wallet.WalletTransaction, int, error) {
	return nil, 0, nil
}
func (m *MockWalletService) ListPendingWithdrawals(ctx context.Context, limit, offset int) ([]wallet.WalletTransaction, int, error) {
	return nil, 0, nil
}
func (m *MockWalletService) ApproveWithdrawal(ctx context.Context, txID, adminID uuid.UUID, txHash string) error {
	return nil
}
func (m *MockWalletService) RejectWithdrawal(ctx context.Context, txID, adminID uuid.UUID, reason string) error {
	return nil
}
func (m *MockWalletService) GetDailyLimitInfo(ctx context.Context, userID uuid.UUID) (float64, float64, error) {
	return 0, 500.0, nil
}

func TestRegister(t *testing.T) {
	repo := NewMockUserRepository()
	walletSvc := &MockWalletService{}
	service := NewAuthService(AuthServiceConfig{
		UserRepo:           repo,
		JWTSecret:          "secret",
		JWTExpiry:          time.Hour,
		RefreshTokenExpiry: time.Hour * 24,
	}).WithWalletService(walletSvc)

	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	password := "Password123!"

	// Test successful registration
	user, err := service.Register(ctx, email, username, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}

	if user.Username != username {
		t.Errorf("expected username %s, got %s", username, user.Username)
	}

	if !walletSvc.CreateCalled {
		t.Error("expected wallet creation to be called")
	}

	if walletSvc.CreatedFor != user.UserID {
		t.Errorf("expected wallet created for user %s, got %s", user.UserID, walletSvc.CreatedFor)
	}

	// Test duplicate registration
	_, err = service.Register(ctx, email, username, password)
	if err == nil {
		t.Error("expected error for duplicate registration, got nil")
	}
}

func TestLogin(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewAuthService(AuthServiceConfig{
		UserRepo:           repo,
		JWTSecret:          "secret",
		JWTExpiry:          time.Hour,
		RefreshTokenExpiry: time.Hour * 24,
	})

	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	password := "Password123!"

	// Register user first
	_, _ = service.Register(ctx, email, username, password)

	// Test successful login
	access, refresh, user, err := service.Login(ctx, email, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if access == "" || refresh == "" {
		t.Error("expected non-empty tokens")
	}

	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}

	// Test invalid password
	_, _, _, err = service.Login(ctx, email, "wrongpassword")
	if err == nil {
		t.Error("expected error for invalid password, got nil")
	}
}
