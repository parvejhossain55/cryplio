package wallet

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TransactionType represents wallet transaction types
type TransactionType string

const (
	TransactionTypeDeposit       TransactionType = "deposit"
	TransactionTypeWithdrawal    TransactionType = "withdrawal"
	TransactionTypeTradeSale     TransactionType = "trade_sale"
	TransactionTypeTradePurchase TransactionType = "trade_purchase"
	TransactionTypeFee           TransactionType = "fee"
	TransactionTypeRefund        TransactionType = "refund"
	TransactionTypeEscrowLock    TransactionType = "escrow_lock"
	TransactionTypeEscrowRelease TransactionType = "escrow_release"
	TransactionTypeDisputeHold   TransactionType = "dispute_hold"
	TransactionTypeDisputeRefund TransactionType = "dispute_refund"
)

// TransactionStatus represents transaction status
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// Wallet represents a user's cryptocurrency wallet
type Wallet struct {
	WalletID      uuid.UUID `db:"wallet_id" json:"wallet_id"`
	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	CryptoID      int       `db:"crypto_id" json:"crypto_id"`
	Address       string    `db:"address" json:"address"`
	Balance       float64   `db:"balance" json:"balance"`
	LockedBalance float64   `db:"locked_balance" json:"locked_balance"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	IsPrimary     bool      `db:"is_primary" json:"is_primary"`
	LastUpdated   time.Time `db:"last_updated" json:"last_updated"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// AvailableBalance returns the available (unlocked) balance
func (w *Wallet) AvailableBalance() float64 {
	return w.Balance - w.LockedBalance
}

// CanDebit checks if the wallet can debit the given amount
func (w *Wallet) CanDebit(amount float64) bool {
	return w.AvailableBalance() >= amount && w.IsActive
}

// Credit adds amount to the wallet balance
func (w *Wallet) Credit(amount float64) {
	w.Balance += amount
	w.LastUpdated = time.Now()
}

// Debit subtracts amount from the wallet balance
func (w *Wallet) Debit(amount float64) error {
	if !w.CanDebit(amount) {
		return fmt.Errorf("insufficient balance")
	}
	w.Balance -= amount
	w.LastUpdated = time.Now()
	return nil
}

// Lock locks the specified amount in the wallet
func (w *Wallet) Lock(amount float64) error {
	if !w.CanDebit(amount) {
		return fmt.Errorf("cannot lock amount: insufficient available balance")
	}
	w.LockedBalance += amount
	w.LastUpdated = time.Now()
	return nil
}

// Unlock releases the locked amount
func (w *Wallet) Unlock(amount float64) {
	w.LockedBalance -= amount
	if w.LockedBalance < 0 {
		w.LockedBalance = 0
	}
	w.LastUpdated = time.Now()
}

// WalletTransaction represents a wallet transaction ledger entry
type WalletTransaction struct {
	TxID               uuid.UUID         `db:"tx_id" json:"tx_id"`
	WalletID           uuid.UUID         `db:"wallet_id" json:"wallet_id"`
	UserID             uuid.UUID         `db:"user_id" json:"user_id"`
	Type               TransactionType   `db:"type" json:"type"`
	Status             TransactionStatus `db:"status" json:"status"`
	Amount             float64           `db:"amount" json:"amount"`
	Fee                float64           `db:"fee" json:"fee"`
	NetAmount          float64           `db:"net_amount" json:"net_amount"`
	BalanceAfter       float64           `db:"balance_after" json:"balance_after"`
	ReferenceID        *uuid.UUID        `db:"reference_id" json:"reference_id,omitempty"` // Links to trade, withdrawal, etc.
	TxHash             *string           `db:"tx_hash" json:"tx_hash,omitempty"`           // Blockchain transaction hash
	Confirmations      int               `db:"confirmations" json:"confirmations"`
	Memo               *string           `db:"memo" json:"memo,omitempty"`
	RequiresApproval   bool              `db:"requires_approval" json:"requires_approval"`
	ApprovedBy         *uuid.UUID        `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt         *time.Time        `db:"approved_at" json:"approved_at,omitempty"`
	DestinationAddress *string           `db:"destination_address" json:"destination_address,omitempty"`
	CreatedAt          time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time         `db:"updated_at" json:"updated_at"`
}

// IsCompleted checks if the transaction is completed
func (t *WalletTransaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsFailed checks if the transaction failed
func (t *WalletTransaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed || t.Status == TransactionStatusCancelled
}

// IsPending checks if the transaction is still pending
func (t *WalletTransaction) IsPending() bool {
	return t.Status == TransactionStatusPending || t.Status == TransactionStatusConfirmed
}

// IsDeposit checks if the transaction is a deposit
func (t *WalletTransaction) IsDeposit() bool {
	return t.Type == TransactionTypeDeposit || t.Type == TransactionTypeRefund || t.Type == TransactionTypeEscrowRelease
}

// IsWithdrawal checks if the transaction is a withdrawal
func (t *WalletTransaction) IsWithdrawal() bool {
	return t.Type == TransactionTypeWithdrawal || t.Type == TransactionTypeTradePurchase || t.Type == TransactionTypeFee
}
