package trade

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"cryplio/pkg/database"

	"github.com/google/uuid"
)

// TradeStatus represents the status of a trade
type TradeStatus string

const (
	TradeStatusPending   TradeStatus = "pending"
	TradeStatusActive    TradeStatus = "active"
	TradeStatusPaid      TradeStatus = "paid"
	TradeStatusReleased  TradeStatus = "released"
	TradeStatusCompleted TradeStatus = "completed"
	TradeStatusCancelled TradeStatus = "cancelled"
	TradeStatusDisputed  TradeStatus = "disputed"
)

// TradeStatusService manages trade status transitions and timers
type TradeStatusService struct {
	db             *database.DB
	paymentWindow  time.Duration
	disputeWindow  time.Duration
	autoExpire     bool
}

// NewTradeStatusService creates a new trade status service
func NewTradeStatusService(db *database.DB, paymentWindow, disputeWindow time.Duration, autoExpire bool) *TradeStatusService {
	return &TradeStatusService{
		db:            db,
		paymentWindow: paymentWindow,
		disputeWindow: disputeWindow,
		autoExpire:    autoExpire,
	}
}

// Trade represents a trade with status information
type Trade struct {
	TradeID           string    `json:"trade_id"`
	AdID              string    `json:"ad_id"`
	BuyerID           string    `json:"buyer_id"`
	SellerID          string    `json:"seller_id"`
	CryptoAmount      float64   `json:"crypto_amount"`
	FiatAmount        float64   `json:"fiat_amount"`
	CryptoSymbol      string    `json:"crypto_symbol"`
	FiatSymbol        string    `json:"fiat_symbol"`
	Rate              float64   `json:"rate"`
	Status            TradeStatus `json:"status"`
	Type              string    `json:"type"`
	PaymentMethod     string    `json:"payment_method"`
	PaymentDetails    string    `json:"payment_details"`
	EscrowID          string    `json:"escrow_id"`
	BlockchainTxHash  string    `json:"blockchain_tx_hash"`
	PaymentWindowMins int       `json:"payment_window_minutes"`
	TimerExpiresAt    *time.Time `json:"timer_expires_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// StatusTransition represents a status transition rule
type StatusTransition struct {
	From      TradeStatus
	To        TradeStatus
	Condition string
	Action    string
}

// GetValidTransitions returns valid status transitions for a trade
func (s *TradeStatusService) GetValidTransitions(currentStatus TradeStatus) []TradeStatus {
	transitions := map[TradeStatus][]TradeStatus{
		TradeStatusPending:   {TradeStatusActive, TradeStatusCancelled},
		TradeStatusActive:    {TradeStatusPaid, TradeStatusCancelled, TradeStatusDisputed},
		TradeStatusPaid:      {TradeStatusReleased, TradeStatusDisputed},
		TradeStatusReleased:  {TradeStatusCompleted},
		TradeStatusCompleted: {}, // Terminal state
		TradeStatusCancelled: {}, // Terminal state
		TradeStatusDisputed:  {TradeStatusCompleted, TradeStatusCancelled}, // Admin resolution
	}

	return transitions[currentStatus]
}

// CanTransition checks if a status transition is valid
func (s *TradeStatusService) CanTransition(from, to TradeStatus) bool {
	validTransitions := s.GetValidTransitions(from)
	for _, valid := range validTransitions {
		if valid == to {
			return true
		}
	}
	return false
}

// UpdateTradeStatus updates a trade's status with validation
func (s *TradeStatusService) UpdateTradeStatus(ctx context.Context, tradeID string, newStatus TradeStatus, notes string) error {
	// Get current trade
	trade, err := s.GetTrade(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("get trade: %w", err)
	}

	// Validate transition
	if !s.CanTransition(trade.Status, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", trade.Status, newStatus)
	}

	// Update status
	now := time.Now().UTC()
	query := `
		UPDATE trades
		SET status = $2, updated_at = $3
		WHERE trade_id = $1
	`

	_, err = s.db.ExecContext(ctx, query, tradeID, newStatus, now)
	if err != nil {
		return fmt.Errorf("update trade status: %w", err)
	}

	// Log status change
	err = s.logStatusChange(ctx, tradeID, trade.Status, newStatus, notes)
	if err != nil {
		log.Printf("Warning: failed to log status change: %v", err)
	}

	// Handle status-specific actions
	err = s.handleStatusChange(ctx, trade, newStatus)
	if err != nil {
		log.Printf("Warning: failed to handle status change: %v", err)
	}

	return nil
}

// StartPaymentTimer starts the payment timer for a trade
func (s *TradeStatusService) StartPaymentTimer(ctx context.Context, tradeID string, windowMinutes int) error {
	expiresAt := time.Now().UTC().Add(time.Duration(windowMinutes) * time.Minute)

	query := `
		UPDATE trades
		SET timer_expires_at = $2, updated_at = $3
		WHERE trade_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, tradeID, expiresAt, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("start payment timer: %w", err)
	}

	// Schedule expiration check
	if s.autoExpire {
		go s.scheduleExpirationCheck(tradeID, expiresAt)
	}

	return nil
}

// GetTrade retrieves a trade by ID
func (s *TradeStatusService) GetTrade(ctx context.Context, tradeID string) (*Trade, error) {
	query := `
		SELECT trade_id, ad_id, buyer_id, seller_id, crypto_amount, fiat_amount, crypto_symbol,
			   fiat_symbol, rate, status, type, payment_method, payment_details, escrow_id,
			   blockchain_tx_hash, payment_window_minutes, timer_expires_at, created_at, updated_at
		FROM trades
		WHERE trade_id = $1
	`

	var trade Trade
	err := s.db.QueryRowContext(ctx, query, tradeID).Scan(
		&trade.TradeID, &trade.AdID, &trade.BuyerID, &trade.SellerID,
		&trade.CryptoAmount, &trade.FiatAmount, &trade.CryptoSymbol,
		&trade.FiatSymbol, &trade.Rate, &trade.Status, &trade.Type,
		&trade.PaymentMethod, &trade.PaymentDetails, &trade.EscrowID,
		&trade.BlockchainTxHash, &trade.PaymentWindowMins, &trade.TimerExpiresAt,
		&trade.CreatedAt, &trade.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade not found")
		}
		return nil, fmt.Errorf("get trade: %w", err)
	}

	return &trade, nil
}

// GetExpiredTrades retrieves trades with expired payment timers
func (s *TradeStatusService) GetExpiredTrades(ctx context.Context) ([]Trade, error) {
	query := `
		SELECT trade_id, ad_id, buyer_id, seller_id, crypto_amount, fiat_amount, crypto_symbol,
			   fiat_symbol, rate, status, type, payment_method, payment_details, escrow_id,
			   blockchain_tx_hash, payment_window_minutes, timer_expires_at, created_at, updated_at
		FROM trades
		WHERE timer_expires_at IS NOT NULL
		  AND timer_expires_at <= $1
		  AND status IN ('active', 'paid')
		ORDER BY timer_expires_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("query expired trades: %w", err)
	}
	defer rows.Close()

	var trades []Trade
	for rows.Next() {
		var trade Trade
		err := rows.Scan(
			&trade.TradeID, &trade.AdID, &trade.BuyerID, &trade.SellerID,
			&trade.CryptoAmount, &trade.FiatAmount, &trade.CryptoSymbol,
			&trade.FiatSymbol, &trade.Rate, &trade.Status, &trade.Type,
			&trade.PaymentMethod, &trade.PaymentDetails, &trade.EscrowID,
			&trade.BlockchainTxHash, &trade.PaymentWindowMins, &trade.TimerExpiresAt,
			&trade.CreatedAt, &trade.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan expired trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, nil
}

// ProcessExpiredTrades processes all expired trades
func (s *TradeStatusService) ProcessExpiredTrades(ctx context.Context) error {
	trades, err := s.GetExpiredTrades(ctx)
	if err != nil {
		return fmt.Errorf("get expired trades: %w", err)
	}

	for _, trade := range trades {
		err := s.processExpiredTrade(ctx, &trade)
		if err != nil {
			log.Printf("Failed to process expired trade %s: %v", trade.TradeID, err)
		}
	}

	return nil
}

// processExpiredTrade processes a single expired trade
func (s *TradeStatusService) processExpiredTrade(ctx context.Context, trade *Trade) error {
	switch trade.Status {
	case TradeStatusActive:
		// Payment not received - cancel trade
		return s.UpdateTradeStatus(ctx, trade.TradeID, TradeStatusCancelled, "Payment window expired")
	case TradeStatusPaid:
		// Payment received but not released - mark for dispute
		return s.UpdateTradeStatus(ctx, trade.TradeID, TradeStatusDisputed, "Release window expired - auto dispute")
	default:
		return fmt.Errorf("unexpected trade status for expiration: %s", trade.Status)
	}
}

// GetTradeTimer returns the remaining time on a trade's payment timer
func (s *TradeStatusService) GetTradeTimer(ctx context.Context, tradeID string) (*time.Duration, string, error) {
	trade, err := s.GetTrade(ctx, tradeID)
	if err != nil {
		return nil, "", fmt.Errorf("get trade: %w", err)
	}

	if trade.TimerExpiresAt == nil {
		return nil, "", fmt.Errorf("no timer set for trade")
	}

	now := time.Now().UTC()
	remaining := trade.TimerExpiresAt.Sub(now)

	if remaining <= 0 {
		return nil, "expired", nil
	}

	return &remaining, "active", nil
}

// ExtendPaymentTimer extends the payment timer for a trade
func (s *TradeStatusService) ExtendPaymentTimer(ctx context.Context, tradeID string, additionalMinutes int) error {
	trade, err := s.GetTrade(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("get trade: %w", err)
	}

	if trade.TimerExpiresAt == nil {
		return fmt.Errorf("no timer set for trade")
	}

	newExpiresAt := trade.TimerExpiresAt.Add(time.Duration(additionalMinutes) * time.Minute)

	query := `
		UPDATE trades
		SET timer_expires_at = $2, updated_at = $3
		WHERE trade_id = $1
	`

	_, err = s.db.ExecContext(ctx, query, tradeID, newExpiresAt, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("extend payment timer: %w", err)
	}

	// Log extension
	err = s.logStatusChange(ctx, tradeID, trade.Status, trade.Status, fmt.Sprintf("Timer extended by %d minutes", additionalMinutes))
	if err != nil {
		log.Printf("Warning: failed to log timer extension: %v", err)
	}

	return nil
}

// CancelTrade cancels a trade with proper validation
func (s *TradeStatusService) CancelTrade(ctx context.Context, tradeID string, reason string, userID string) error {
	trade, err := s.GetTrade(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("get trade: %w", err)
	}

	// Check if user can cancel this trade
	if !s.canUserCancelTrade(trade, userID) {
		return fmt.Errorf("user not authorized to cancel this trade")
	}

	// Check if trade can be cancelled
	if !s.CanTransition(trade.Status, TradeStatusCancelled) {
		return fmt.Errorf("trade cannot be cancelled in current status: %s", trade.Status)
	}

	return s.UpdateTradeStatus(ctx, tradeID, TradeStatusCancelled, fmt.Sprintf("Cancelled by %s: %s", userID, reason))
}

// canUserCancelTrade checks if a user can cancel a trade
func (s *TradeStatusService) canUserCancelTrade(trade *Trade, userID string) bool {
	// Buyer can cancel if payment not yet made
	if userID == trade.BuyerID && trade.Status == TradeStatusActive {
		return true
	}

	// Seller can cancel if trade is pending
	if userID == trade.SellerID && trade.Status == TradeStatusPending {
		return true
	}

	// Admin can cancel any trade (handled elsewhere)
	return false
}

// handleStatusChange performs actions based on status changes
func (s *TradeStatusService) handleStatusChange(ctx context.Context, trade *Trade, newStatus TradeStatus) error {
	switch newStatus {
	case TradeStatusActive:
		// Start payment timer when trade becomes active
		if trade.PaymentWindowMins > 0 {
			return s.StartPaymentTimer(ctx, trade.TradeID, trade.PaymentWindowMins)
		}
	case TradeStatusCompleted, TradeStatusCancelled:
		// Clear timer when trade is completed or cancelled
		return s.clearPaymentTimer(ctx, trade.TradeID)
	}

	return nil
}

// clearPaymentTimer removes the payment timer from a trade
func (s *TradeStatusService) clearPaymentTimer(ctx context.Context, tradeID string) error {
	query := `
		UPDATE trades
		SET timer_expires_at = NULL, updated_at = $2
		WHERE trade_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, tradeID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("clear payment timer: %w", err)
	}

	return nil
}

// logStatusChange logs a status change to the audit log
func (s *TradeStatusService) logStatusChange(ctx context.Context, tradeID string, fromStatus, toStatus TradeStatus, notes string) error {
	query := `
		INSERT INTO trade_status_log (log_id, trade_id, from_status, to_status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query,
		uuid.New().String(), tradeID, fromStatus, toStatus, notes, time.Now().UTC(),
	)

	if err != nil {
		return fmt.Errorf("log status change: %w", err)
	}

	return nil
}

// scheduleExpirationCheck schedules a check for trade expiration
func (s *TradeStatusService) scheduleExpirationCheck(tradeID string, expiresAt time.Time) {
	delay := time.Until(expiresAt)
	if delay <= 0 {
		return
	}

	time.AfterFunc(delay, func() {
		ctx := context.Background()
		err := s.ProcessExpiredTrades(ctx)
		if err != nil {
			log.Printf("Failed to process expired trades: %v", err)
		}
	})
}

// GetTradeStatusHistory returns the status change history for a trade
func (s *TradeStatusService) GetTradeStatusHistory(ctx context.Context, tradeID string) ([]StatusHistory, error) {
	query := `
		SELECT log_id, trade_id, from_status, to_status, notes, created_at
		FROM trade_status_log
		WHERE trade_id = $1
		ORDER BY created_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tradeID)
	if err != nil {
		return nil, fmt.Errorf("query status history: %w", err)
	}
	defer rows.Close()

	var history []StatusHistory
	for rows.Next() {
		var h StatusHistory
		err := rows.Scan(
			&h.LogID, &h.TradeID, &h.FromStatus, &h.ToStatus, &h.Notes, &h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan status history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}

// StatusHistory represents a status change in the trade's history
type StatusHistory struct {
	LogID      string     `json:"log_id"`
	TradeID    string     `json:"trade_id"`
	FromStatus TradeStatus `json:"from_status"`
	ToStatus   TradeStatus `json:"to_status"`
	Notes      string     `json:"notes"`
	CreatedAt  time.Time  `json:"created_at"`
}

// StartWorker starts the background worker for processing expired trades
func (s *TradeStatusService) StartWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx := context.Background()
			err := s.ProcessExpiredTrades(ctx)
			if err != nil {
				log.Printf("Failed to process expired trades: %v", err)
			}
		}
	}
}
