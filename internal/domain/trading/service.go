package trading

import (
	"context"
	"cryplio/internal/domain/dispute"
	"cryplio/internal/domain/identity"
	"cryplio/internal/domain/notification"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TradeService interface {
	// Ads
	ListActiveAds(ctx context.Context, filter AdFilter) ([]TradeAd, int, error)
	GetAd(ctx context.Context, id uuid.UUID) (*TradeAd, error)
	CreateAd(ctx context.Context, ad *TradeAd) error
	UpdateAd(ctx context.Context, adID, userID uuid.UUID, updates *TradeAd) error
	DeleteAd(ctx context.Context, adID, userID uuid.UUID) error

	// Trades
	InitiateTrade(ctx context.Context, adID, buyerID uuid.UUID, amount float64) (*Trade, error)
	ListTrades(ctx context.Context, userID uuid.UUID, role string) ([]Trade, error)
	ListAllTrades(ctx context.Context, status string) ([]Trade, error)
	CountTrades(ctx context.Context, status string) (int, error)
	GetTrade(ctx context.Context, id uuid.UUID) (*Trade, error)
	MarkAsPaid(ctx context.Context, tradeID, userID uuid.UUID) error
	ReleaseEscrow(ctx context.Context, tradeID, userID uuid.UUID) error
	CancelTrade(ctx context.Context, tradeID, userID uuid.UUID, reason string) error
	DisputeTrade(ctx context.Context, tradeID, userID uuid.UUID, reasonCode string, reasonText string) (*Trade, error)
	ReconcileExpiredTrades(ctx context.Context) (int, error)
	FlagAutoDisputesForOverduePaidTrades(ctx context.Context, gracePeriod time.Duration) (int, error)

	// Messages
	SendMessage(ctx context.Context, tradeID, senderID uuid.UUID, content string) (*TradeMessage, error)
	SendFileMessage(ctx context.Context, tradeID, senderID uuid.UUID, fileURL, mimeType string, fileSize int) (*TradeMessage, error)
	GetChatHistory(ctx context.Context, tradeID, userID uuid.UUID) ([]TradeMessage, error)

	// Feedback
	LeaveFeedback(ctx context.Context, tradeID, userID uuid.UUID, rating FeedbackRating, comment string) error

	// Merchant Management
	ListUserAds(ctx context.Context, userID uuid.UUID) ([]TradeAd, int, error)
	ToggleAdStatus(ctx context.Context, adID, userID uuid.UUID) error
}

type tradeService struct {
	tradeRepo           TradeRepository
	identityRepo        identity.UserRepository
	disputeRepo         dispute.Repository
	escrowClient        EscrowContractClient
	notificationService notification.Service
}

func NewTradeService(tradeRepo TradeRepository, identityRepo identity.UserRepository, disputeRepo dispute.Repository, escrowClient EscrowContractClient, notificationService notification.Service) TradeService {
	return &tradeService{
		tradeRepo:           tradeRepo,
		identityRepo:        identityRepo,
		disputeRepo:         disputeRepo,
		escrowClient:        escrowClient,
		notificationService: notificationService,
	}
}

func (s *tradeService) ListActiveAds(ctx context.Context, filter AdFilter) ([]TradeAd, int, error) {
	status := TradeAdStatusActive
	filter.Status = &status
	return s.tradeRepo.ListAds(ctx, filter)
}

func (s *tradeService) GetAd(ctx context.Context, id uuid.UUID) (*TradeAd, error) {
	return s.tradeRepo.GetAdByID(ctx, id)
}

func (s *tradeService) CreateAd(ctx context.Context, ad *TradeAd) error {
	if ad.AdID == uuid.Nil {
		ad.AdID = uuid.New()
	}
	return s.tradeRepo.CreateAd(ctx, ad)
}

func (s *tradeService) UpdateAd(ctx context.Context, adID, userID uuid.UUID, updates *TradeAd) error {
	ad, err := s.tradeRepo.GetAdByID(ctx, adID)
	if err != nil {
		return err
	}
	if ad == nil {
		return errors.New("ad not found")
	}
	if ad.UserID != userID {
		return errors.New("unauthorized")
	}

	if updates.Type != "" {
		ad.Type = updates.Type
	}
	if updates.CryptoID != 0 {
		ad.CryptoID = updates.CryptoID
	}
	if updates.FiatID != 0 {
		ad.FiatID = updates.FiatID
	}
	if updates.PriceType != "" {
		ad.PriceType = updates.PriceType
	}
	if updates.Price > 0 {
		ad.Price = updates.Price
	}
	if updates.FloatingMarkup != nil {
		ad.FloatingMarkup = updates.FloatingMarkup
	}
	if updates.MinAmount > 0 {
		ad.MinAmount = updates.MinAmount
	}
	if updates.MaxAmount > 0 {
		ad.MaxAmount = updates.MaxAmount
	}
	if len(updates.PaymentMethods) > 0 {
		ad.PaymentMethods = updates.PaymentMethods
	}
	if updates.TradeTerms != nil {
		ad.TradeTerms = updates.TradeTerms
	}
	if updates.PaymentWindowMinutes > 0 {
		ad.PaymentWindowMinutes = updates.PaymentWindowMinutes
	}
	if updates.Timezone != "" {
		ad.Timezone = updates.Timezone
	}

	return s.tradeRepo.UpdateAd(ctx, ad)
}

func (s *tradeService) DeleteAd(ctx context.Context, adID, userID uuid.UUID) error {
	ad, err := s.tradeRepo.GetAdByID(ctx, adID)
	if err != nil {
		return err
	}
	if ad == nil {
		return errors.New("ad not found")
	}
	if ad.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.tradeRepo.DeleteAd(ctx, adID)
}

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

	// 1. Fetch Ad
	ad, err := s.tradeRepo.GetAdByID(ctx, adID)
	if err != nil {
		return nil, fmt.Errorf("get ad: %w", err)
	}
	if ad == nil {
		return nil, errors.New("advertisement not found")
	}

	// 2. Validate Ad Status
	if !ad.IsActive() {
		return nil, errors.New("advertisement is not active")
	}

	// 3. Prevent Self-Trading
	if ad.UserID == buyerID {
		return nil, errors.New("you cannot trade with your own advertisement")
	}

	// 4. Validate Amount
	if amount < ad.MinAmount || amount > ad.MaxAmount {
		return nil, fmt.Errorf("amount must be between %.2f and %.2f", ad.MinAmount, ad.MaxAmount)
	}

	// 7. Create Trade
	trade := &Trade{
		TradeID:              uuid.New(),
		AdID:                 ad.AdID,
		BuyerID:              buyerID,
		SellerID:             ad.UserID,
		CryptoID:             ad.CryptoID,
		FiatID:               ad.FiatID,
		CryptoAmount:         amount / ad.Price, // Simple calculation for now
		FiatAmount:           amount,
		ExchangeRate:         ad.Price,
		PaymentMethod:        ad.PaymentMethods[0], // Use first method for now
		PriceType:            ad.PriceType,
		AgreedPrice:          ad.Price,
		Status:               TradeStatusPending,
		PaymentWindowMinutes: ad.PaymentWindowMinutes,
	}

	// 8. Lock Escrow on Blockchain
	txHash, contractAddr, err := s.escrowClient.Lock(ctx, trade)
	if err != nil {
		return nil, fmt.Errorf("blockchain escrow lock failed: %w", err)
	}
	now := time.Now()
	trade.EscrowTxnHash = &txHash
	trade.EscrowContractAddress = &contractAddr
	trade.Status = TradeStatusActive // Active once locked
	trade.StartedAt = &now

	err = s.tradeRepo.CreateTrade(ctx, trade)
	if err != nil {
		return nil, fmt.Errorf("save trade: %w", err)
	}

	// Notify both parties
	if s.notificationService != nil {
		seller, _ := s.identityRepo.GetByID(ctx, ad.UserID)
		if seller != nil {
			msg := fmt.Sprintf("New trade initiated for your ad. Trade ID: %s", trade.TradeID.String())
			_ = s.notificationService.Notify(ctx, ad.UserID, notification.NotificationTypeTradeStarted, "New Trade Started", msg, nil)
		}
		buyer, _ := s.identityRepo.GetByID(ctx, buyerID)
		if buyer != nil {
			msg := fmt.Sprintf("You have started a new trade. Trade ID: %s", trade.TradeID.String())
			_ = s.notificationService.Notify(ctx, buyerID, notification.NotificationTypeTradeStarted, "Trade Started", msg, nil)
		}
	}

	return trade, nil
}

func (s *tradeService) ListUserAds(ctx context.Context, userID uuid.UUID) ([]TradeAd, int, error) {
	filter := AdFilter{
		UserID: &userID,
		Limit:  100, // Show all for now
	}
	return s.tradeRepo.ListAds(ctx, filter)
}

func (s *tradeService) ToggleAdStatus(ctx context.Context, adID, userID uuid.UUID) error {
	ad, err := s.tradeRepo.GetAdByID(ctx, adID)
	if err != nil {
		return err
	}
	if ad == nil {
		return errors.New("ad not found")
	}
	if ad.UserID != userID {
		return errors.New("unauthorized")
	}

	ad.IsPaused = !ad.IsPaused
	return s.tradeRepo.UpdateAd(ctx, ad)
}

func (s *tradeService) ListTrades(ctx context.Context, userID uuid.UUID, role string) ([]Trade, error) {
	return s.tradeRepo.ListTrades(ctx, userID, role)
}

func (s *tradeService) ListAllTrades(ctx context.Context, status string) ([]Trade, error) {
	return s.tradeRepo.ListAllTrades(ctx, status)
}

func (s *tradeService) CountTrades(ctx context.Context, status string) (int, error) {
	return s.tradeRepo.CountTrades(ctx, status)
}

func (s *tradeService) GetTrade(ctx context.Context, id uuid.UUID) (*Trade, error) {
	return s.tradeRepo.GetTradeByID(ctx, id)
}

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
	// Notify seller that buyer marked as paid
	if s.notificationService != nil {
		msg := fmt.Sprintf("Buyer marked trade %s as paid. Please verify and release escrow.", trade.TradeID.String())
		_ = s.notificationService.Notify(ctx, trade.SellerID, notification.NotificationTypeTradePaid, "Trade Marked as Paid", msg, nil)
	}
	return nil
}

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

	// 2. Release Escrow on Blockchain
	txHash, err := s.escrowClient.Release(ctx, trade)
	if err != nil {
		return fmt.Errorf("blockchain escrow release failed: %w", err)
	}
	trade.EscrowTxnHash = &txHash

	// 3. Auto-complete trade after escrow release
	trade.Complete()

	if err := s.tradeRepo.UpdateTrade(ctx, trade); err != nil {
		return err
	}
	// Notify both parties that trade is completed
	if s.notificationService != nil {
		msg := fmt.Sprintf("Trade %s completed successfully. Escrow has been released.", trade.TradeID.String())
		_ = s.notificationService.Notify(ctx, trade.BuyerID, notification.NotificationTypeTradeCompleted, "Trade Completed", msg, nil)
		_ = s.notificationService.Notify(ctx, trade.SellerID, notification.NotificationTypeTradeCompleted, "Trade Completed", msg, nil)
	}
	return nil
}

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
	// Notify the other party
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

	// Create Dispute
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

	// Update Trade
	trade.Status = TradeStatusDisputed
	trade.DisputeID = &d.DisputeID
	if err := s.tradeRepo.UpdateTrade(ctx, trade); err != nil {
		return nil, fmt.Errorf("update trade status: %w", err)
	}

	// Notify both parties
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

func (s *tradeService) SendMessage(ctx context.Context, tradeID, senderID uuid.UUID, content string) (*TradeMessage, error) {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.BuyerID != senderID && trade.SellerID != senderID {
		return nil, errors.New("unauthorized")
	}

	msg := &TradeMessage{
		MessageID:   uuid.New(),
		TradeID:     tradeID,
		SenderID:    senderID,
		MessageType: TradeMessageTypeText,
		Content:     &content,
	}

	err = s.tradeRepo.CreateTradeMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *tradeService) SendFileMessage(ctx context.Context, tradeID, senderID uuid.UUID, fileURL, mimeType string, fileSize int) (*TradeMessage, error) {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.BuyerID != senderID && trade.SellerID != senderID {
		return nil, errors.New("unauthorized")
	}

	messageType := TradeMessageTypeFile
	if strings.HasPrefix(mimeType, "image/") {
		messageType = TradeMessageTypeImage
	}

	msg := &TradeMessage{
		MessageID:    uuid.New(),
		TradeID:      tradeID,
		SenderID:     senderID,
		MessageType:  messageType,
		FileURL:      &fileURL,
		FileMimeType: &mimeType,
		FileSize:     &fileSize,
	}

	err = s.tradeRepo.CreateTradeMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *tradeService) GetChatHistory(ctx context.Context, tradeID, userID uuid.UUID) ([]TradeMessage, error) {
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

	return s.tradeRepo.ListTradeMessages(ctx, tradeID)
}

func (s *tradeService) LeaveFeedback(ctx context.Context, tradeID, userID uuid.UUID, rating FeedbackRating, comment string) error {
	// Get trade
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return err
	}
	if trade == nil {
		return errors.New("trade not found")
	}

	// Check if trade is completed
	if trade.Status != TradeStatusCompleted {
		return errors.New("can only leave feedback on completed trades")
	}

	// Check if user participated in trade
	if trade.BuyerID != userID && trade.SellerID != userID {
		return errors.New("unauthorized")
	}

	// Check if feedback already exists
	existing, err := s.tradeRepo.GetFeedbackByTrade(ctx, tradeID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("feedback already exists for this trade")
	}

	// Determine recipient
	var recipientID uuid.UUID
	if trade.BuyerID == userID {
		recipientID = trade.SellerID
	} else {
		recipientID = trade.BuyerID
	}

	// Create feedback
	feedback := &TradeFeedback{
		FeedbackID: uuid.New(),
		TradeID:    tradeID,
		FromUserID: userID,
		ToUserID:   recipientID,
		Rating:     rating,
		Comment:    &comment,
	}

	if err := s.tradeRepo.CreateFeedback(ctx, feedback); err != nil {
		return fmt.Errorf("create feedback: %w", err)
	}

	return nil
}

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

func (s *tradeService) FlagAutoDisputesForOverduePaidTrades(ctx context.Context, gracePeriod time.Duration) (int, error) {
	if gracePeriod <= 0 {
		gracePeriod = time.Hour
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
