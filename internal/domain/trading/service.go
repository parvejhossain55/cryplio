package trading

import (
	"context"
	"cryplio/internal/domain/identity"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type TradeService interface {
	// Ads
	ListActiveAds(ctx context.Context, filter AdFilter) ([]TradeAd, int, error)
	GetAd(ctx context.Context, id uuid.UUID) (*TradeAd, error)
	CreateAd(ctx context.Context, ad *TradeAd) error

	// Trades
	InitiateTrade(ctx context.Context, adID, buyerID uuid.UUID, amount float64) (*Trade, error)
	ListTrades(ctx context.Context, userID uuid.UUID, role string) ([]Trade, error)
	GetTrade(ctx context.Context, id uuid.UUID) (*Trade, error)
	MarkAsPaid(ctx context.Context, tradeID, userID uuid.UUID) error
	ReleaseEscrow(ctx context.Context, tradeID, userID uuid.UUID) error
	CancelTrade(ctx context.Context, tradeID, userID uuid.UUID, reason string) error

	// Messages
	SendMessage(ctx context.Context, tradeID, senderID uuid.UUID, content string) (*TradeMessage, error)
	GetChatHistory(ctx context.Context, tradeID, userID uuid.UUID) ([]TradeMessage, error)

	// Merchant Management
	ListUserAds(ctx context.Context, userID uuid.UUID) ([]TradeAd, int, error)
	ToggleAdStatus(ctx context.Context, adID, userID uuid.UUID) error
}

type tradeService struct {
	tradeRepo    TradeRepository
	identityRepo identity.UserRepository
}

func NewTradeService(tradeRepo TradeRepository, identityRepo identity.UserRepository) TradeService {
	return &tradeService{
		tradeRepo:    tradeRepo,
		identityRepo: identityRepo,
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

func (s *tradeService) InitiateTrade(ctx context.Context, adID, buyerID uuid.UUID, amount float64) (*Trade, error) {
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

	// 5. CRITICAL: Check if Buyer is blocked by Seller (Ad Creator) - FR-136
	isBlocked, err := s.identityRepo.IsBlocked(ctx, ad.UserID, buyerID)
	if err != nil {
		return nil, fmt.Errorf("check block status: %w", err)
	}
	if isBlocked {
		return nil, errors.New("you are blocked by this user and cannot initiate trades")
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

	err = s.tradeRepo.CreateTrade(ctx, trade)
	if err != nil {
		return nil, fmt.Errorf("save trade: %w", err)
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
	return s.tradeRepo.UpdateTrade(ctx, trade)
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
	// In a real app, this would also trigger the actual blockchain release
	return s.tradeRepo.UpdateTrade(ctx, trade)
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
	return s.tradeRepo.UpdateTrade(ctx, trade)
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
