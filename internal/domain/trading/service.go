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

	// 6. Fetch Buyer Info for KYC check
	buyer, err := s.identityRepo.GetByID(ctx, buyerID)
	if err != nil {
		return nil, fmt.Errorf("get buyer: %w", err)
	}
	if buyer.KYCLevel < ad.RequiresKYCLevel {
		return nil, fmt.Errorf("this advertisement requires KYC Level %d", ad.RequiresKYCLevel)
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
