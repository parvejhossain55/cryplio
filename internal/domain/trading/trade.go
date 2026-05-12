package trading

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cryplio/internal/domain/dispute"
	"cryplio/internal/domain/notification"

	"github.com/google/uuid"
)

// InitiateTrade validates preconditions, locks escrow on-chain, and persists a new Trade.
func (s *tradeService) InitiateTrade(ctx context.Context, adID, buyerID uuid.UUID, amount float64) (*Trade, error) {
	buyer, err := s.identityRepo.GetByID(ctx, buyerID)
	if err != nil {
		return nil, fmt.Errorf("get buyer: %w", err)
	}
	if buyer == nil {
		return nil, errors.New("buyer not found")
	}
	if !buyer.EmailVerified {
		return nil, errors.New("email verification is required before initiating trades")
	}

	// 1. Fetch Ad.
	ad, err := s.tradeRepo.GetAdByID(ctx, adID)
	if err != nil {
		return nil, fmt.Errorf("get ad: %w", err)
	}
	if ad == nil {
		return nil, errors.New("advertisement not found")
	}

	// 2. Validate Ad Status.
	if !ad.IsActive() {
		return nil, errors.New("advertisement is not active")
	}

	// 3. Prevent Self-Trading.
	if ad.UserID == buyerID {
		return nil, errors.New("you cannot trade with your own advertisement")
	}

	// 4. Validate Amount.
	if amount < ad.MinAmount || amount > ad.MaxAmount {
		return nil, fmt.Errorf("amount must be between %.2f and %.2f", ad.MinAmount, ad.MaxAmount)
	}

	if len(ad.PaymentMethods) == 0 {
		return nil, errors.New("advertisement has no payment methods configured")
	}

	// 5. Fetch Blockchain Addresses.
	buyerWallet, err := s.walletRepo.GetByUser(ctx, buyerID)
	if err != nil || buyerWallet == nil {
		return nil, fmt.Errorf("buyer wallet not found: %w", err)
	}

	sellerWallet, err := s.walletRepo.GetByUser(ctx, ad.UserID)
	if err != nil || sellerWallet == nil {
		return nil, fmt.Errorf("seller wallet not found: %w", err)
	}

	asset, err := s.platformRepo.GetCryptoAsset(ctx, ad.CryptoID)
	if err != nil || asset == nil {
		return nil, fmt.Errorf("crypto asset not found: %w", err)
	}

	// 6. Build Trade record.
	trade := &Trade{
		TradeID:              uuid.New(),
		AdID:                 ad.AdID,
		BuyerID:              buyerID,
		SellerID:             ad.UserID,
		CryptoID:             ad.CryptoID,
		FiatID:               ad.FiatID,
		CryptoAmount:         amount / ad.Price, // Simple calculation for now.
		FiatAmount:           amount,
		ExchangeRate:         ad.Price,
		PaymentMethod:        ad.PaymentMethods[0], // Use first method for now.
		PriceType:            ad.PriceType,
		AgreedPrice:          ad.Price,
		Status:               TradeStatusPending,
		PaymentWindowMinutes: ad.PaymentWindowMinutes,
		BuyerAddress:         &buyerWallet.Address,
		SellerAddress:        &sellerWallet.Address,
		TokenAddress:         asset.ContractAddress,
	}

	// 7. Lock Escrow on Blockchain.
	txHash, contractAddr, err := s.escrowClient.Lock(ctx, trade)
	if err != nil {
		return nil, fmt.Errorf("blockchain escrow lock failed: %w", err)
	}
	now := time.Now()
	trade.EscrowTxnHash = &txHash
	trade.EscrowContractAddress = &contractAddr
	trade.Status = TradeStatusActive // Active once locked.
	trade.StartedAt = &now

	if err = s.tradeRepo.CreateTrade(ctx, trade); err != nil {
		return nil, fmt.Errorf("save trade: %w", err)
	}

	// Notify both parties.
	if s.notificationService != nil {
		seller, _ := s.identityRepo.GetByID(ctx, ad.UserID)
		if seller != nil {
			msg := fmt.Sprintf("New trade initiated for your ad. Trade ID: %s", trade.TradeID.String())
			_ = s.notificationService.Notify(ctx, ad.UserID, notification.NotificationTypeTradeStarted, "New Trade Started", msg, nil)
		}
		if buyer != nil {
			msg := fmt.Sprintf("You have started a new trade. Trade ID: %s", trade.TradeID.String())
			_ = s.notificationService.Notify(ctx, buyerID, notification.NotificationTypeTradeStarted, "Trade Started", msg, nil)
		}
	}

	return trade, nil
}

// ListTrades returns all trades where userID participated in the given role.
func (s *tradeService) ListTrades(ctx context.Context, userID uuid.UUID, role string) ([]Trade, error) {
	return s.tradeRepo.ListTrades(ctx, userID, role)
}

// ListAllTrades returns every trade matching the optional status filter.
func (s *tradeService) ListAllTrades(ctx context.Context, status string) ([]Trade, error) {
	return s.tradeRepo.ListAllTrades(ctx, status)
}

// CountTrades returns the number of trades that match the optional status filter.
func (s *tradeService) CountTrades(ctx context.Context, status string) (int, error) {
	return s.tradeRepo.CountTrades(ctx, status)
}

// GetTrade returns a single trade by its ID.
func (s *tradeService) GetTrade(ctx context.Context, id uuid.UUID) (*Trade, error) {
	return s.tradeRepo.GetTradeByID(ctx, id)
}

// MarkAsPaid records that the buyer has sent payment and notifies the seller.
func (s *tradeService) MarkAsPaid(ctx context.Context, tradeID, userID uuid.UUID) error {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return err
	}
	if trade == nil {
		return errors.New("trade not found")
	}
	if trade.BuyerID != userID {
		return errors.New("unauthorized")
	}
	if trade.Status != TradeStatusPending && trade.Status != TradeStatusActive {
		return errors.New("invalid trade status")
	}

	trade.MarkAsPaid()
	if err := s.tradeRepo.UpdateTrade(ctx, trade); err != nil {
		return err
	}

	// Notify seller that buyer marked as paid.
	if s.notificationService != nil {
		msg := fmt.Sprintf("Buyer marked trade %s as paid. Please verify and release escrow.", trade.TradeID.String())
		_ = s.notificationService.Notify(ctx, trade.SellerID, notification.NotificationTypeTradePaid, "Trade Marked as Paid", msg, nil)
	}
	return nil
}

// ReleaseEscrow releases on-chain escrow, completes the trade, and notifies both parties.
func (s *tradeService) ReleaseEscrow(ctx context.Context, tradeID, userID uuid.UUID) error {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return err
	}
	if trade == nil {
		return errors.New("trade not found")
	}
	if trade.SellerID != userID {
		return errors.New("unauthorized")
	}
	if trade.Status != TradeStatusPaid {
		return errors.New("trade must be marked as paid before releasing escrow")
	}

	trade.Release()

	// Release Escrow on Blockchain.
	txHash, err := s.escrowClient.Release(ctx, trade)
	if err != nil {
		return fmt.Errorf("blockchain escrow release failed: %w", err)
	}
	trade.EscrowTxnHash = &txHash

	// Auto-complete trade after escrow release.
	trade.Complete()

	if err := s.tradeRepo.UpdateTrade(ctx, trade); err != nil {
		// CRITICAL: Blockchain release succeeded but DB update failed.
		// The trade is now inconsistent. Manual reconciliation required.
		// TODO: Implement outbox/saga pattern to ensure atomicity.
		return fmt.Errorf("CRITICAL: escrow released on-chain (txHash=%s) but DB update failed — manual reconciliation required: %w", txHash, err)
	}

	// Notify both parties that trade is completed.
	if s.notificationService != nil {
		msg := fmt.Sprintf("Trade %s completed successfully. Escrow has been released.", trade.TradeID.String())
		_ = s.notificationService.Notify(ctx, trade.BuyerID, notification.NotificationTypeTradeCompleted, "Trade Completed", msg, nil)
		_ = s.notificationService.Notify(ctx, trade.SellerID, notification.NotificationTypeTradeCompleted, "Trade Completed", msg, nil)
	}
	return nil
}

// CancelTrade cancels a trade and notifies the other party.
func (s *tradeService) CancelTrade(ctx context.Context, tradeID, userID uuid.UUID, reason string) error {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return err
	}
	if trade == nil {
		return errors.New("trade not found")
	}
	if trade.BuyerID != userID && trade.SellerID != userID {
		return errors.New("unauthorized")
	}
	if !trade.CanCancel() {
		return errors.New("trade cannot be cancelled at this stage")
	}

	trade.Cancel(reason)
	if err := s.tradeRepo.UpdateTrade(ctx, trade); err != nil {
		return err
	}

	// Notify the other party.
	if s.notificationService != nil {
		msg := fmt.Sprintf("Trade %s has been cancelled. Reason: %s", trade.TradeID.String(), reason)
		otherID := trade.SellerID
		if trade.SellerID == userID {
			otherID = trade.BuyerID
		}
		_ = s.notificationService.Notify(ctx, otherID, notification.NotificationTypeTradeCancelled, "Trade Cancelled", msg, nil)
	}
	return nil
}

// DisputeTrade opens a dispute record, marks the trade as disputed, and notifies both parties.
func (s *tradeService) DisputeTrade(ctx context.Context, tradeID, userID uuid.UUID, reasonCode string, reasonText string) (*Trade, error) {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.BuyerID != userID && trade.SellerID != userID {
		return nil, errors.New("unauthorized")
	}
	if trade.Status != TradeStatusPaid && trade.Status != TradeStatusActive {
		return nil, errors.New("trade cannot be disputed in its current status")
	}

	// Create Dispute record.
	d := &dispute.Dispute{
		DisputeID:  uuid.New(),
		TradeID:    tradeID,
		RaisedBy:   userID,
		ReasonCode: reasonCode,
		ReasonText: &reasonText,
		Status:     dispute.DisputeStatusPending,
	}

	if err := s.disputeRepo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("create dispute: %w", err)
	}

	// Update Trade status.
	trade.Status = TradeStatusDisputed
	trade.DisputeID = &d.DisputeID
	if err := s.tradeRepo.UpdateTrade(ctx, trade); err != nil {
		return nil, fmt.Errorf("update trade status: %w", err)
	}

	// Notify both parties.
	if s.notificationService != nil {
		msg := fmt.Sprintf("A dispute has been raised on trade %s. Reason: %s", trade.TradeID.String(), reasonText)
		otherID := trade.SellerID
		if trade.SellerID == userID {
			otherID = trade.BuyerID
		}
		_ = s.notificationService.Notify(ctx, userID, notification.NotificationTypeTradeDisputed, "Trade Disputed", msg, nil)
		_ = s.notificationService.Notify(ctx, otherID, notification.NotificationTypeTradeDisputed, "Trade Disputed", msg, nil)
	}

	return trade, nil
}

// ReconcileExpiredTrades marks pending/active trades whose payment window has passed as expired.
func (s *tradeService) ReconcileExpiredTrades(ctx context.Context) (int, error) {
	now := time.Now()
	trades, err := s.tradeRepo.ListExpiredPendingTrades(ctx, now)
	if err != nil {
		return 0, err
	}

	updated := 0
	for i := range trades {
		trade := trades[i]
		if trade.Status != TradeStatusPending && trade.Status != TradeStatusActive {
			continue
		}

		trade.Status = TradeStatusExpired
		trade.ExpiredAt = &now
		reason := "auto-cancelled: payment window expired"
		trade.CancelReason = &reason

		if err := s.tradeRepo.UpdateTrade(ctx, &trade); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

// FlagAutoDisputesForOverduePaidTrades sets IsAutoDisputeTriggered on trades that have been
// in the Paid state longer than the provided gracePeriod. Uses configured default if not specified.
func (s *tradeService) FlagAutoDisputesForOverduePaidTrades(ctx context.Context, gracePeriod time.Duration) (int, error) {
	if gracePeriod <= 0 {
		// Use configured grace period or default to 1 hour
		if s.cfg != nil && s.cfg.TradeAutoDisputeGracePeriod > 0 {
			gracePeriod = s.cfg.TradeAutoDisputeGracePeriod
		} else {
			gracePeriod = time.Hour
		}
	}

	threshold := time.Now().Add(-gracePeriod)
	trades, err := s.tradeRepo.ListPaidTradesPastGrace(ctx, threshold)
	if err != nil {
		return 0, err
	}

	updated := 0
	for i := range trades {
		trade := trades[i]
		if trade.IsAutoDisputeTriggered {
			continue
		}

		trade.IsAutoDisputeTriggered = true
		if err := s.tradeRepo.UpdateTrade(ctx, &trade); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}
