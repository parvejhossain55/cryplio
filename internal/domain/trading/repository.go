package trading

import (
	"context"

	"github.com/google/uuid"
)

// TradeRepository defines the persistence interface for trading and ads.
type TradeRepository interface {
	// Trade Ads
	CreateAd(ctx context.Context, ad *TradeAd) error
	GetAdByID(ctx context.Context, id uuid.UUID) (*TradeAd, error)
	ListAds(ctx context.Context, filter AdFilter) ([]TradeAd, int, error)
	UpdateAd(ctx context.Context, ad *TradeAd) error
	DeleteAd(ctx context.Context, id uuid.UUID) error

	// Trades
	CreateTrade(ctx context.Context, trade *Trade) error
	GetTradeByID(ctx context.Context, id uuid.UUID) (*Trade, error)
	ListTrades(ctx context.Context, userID uuid.UUID, role string) ([]Trade, error)
	UpdateTrade(ctx context.Context, trade *Trade) error

	// Messages
	CreateTradeMessage(ctx context.Context, msg *TradeMessage) error
	ListTradeMessages(ctx context.Context, tradeID uuid.UUID) ([]TradeMessage, error)
}

// AdFilter holds filtering options for trade ads
type AdFilter struct {
	Type           *AdType
	CryptoID       *int
	FiatID         *int
	PaymentMethods []int
	MinAmount      *float64
	UserID         *uuid.UUID
	Status         *TradeAdStatus
	Limit          int
	Offset         int
}
