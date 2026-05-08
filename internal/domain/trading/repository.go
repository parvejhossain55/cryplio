package trading

import (
	"context"
	"time"

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
	ListAllTrades(ctx context.Context, status string) ([]Trade, error)
	CountTrades(ctx context.Context, status string) (int, error)
	ListExpiredPendingTrades(ctx context.Context, now time.Time) ([]Trade, error)
	ListPaidTradesPastGrace(ctx context.Context, threshold time.Time) ([]Trade, error)
	UpdateTrade(ctx context.Context, trade *Trade) error

	// Messages
	CreateTradeMessage(ctx context.Context, msg *TradeMessage) error
	ListTradeMessages(ctx context.Context, tradeID uuid.UUID) ([]TradeMessage, error)

	// Feedback
	CreateFeedback(ctx context.Context, feedback *TradeFeedback) error
	GetFeedbackByTrade(ctx context.Context, tradeID uuid.UUID) (*TradeFeedback, error)
}

// AdFilter holds filtering options for trade ads
type AdFilter struct {
	Type           *AdType
	CryptoID       *int
	FiatID         *int
	FiatCode       *string
	PaymentMethods []int
	MinAmount      *float64
	UserID         *uuid.UUID
	Status         *TradeAdStatus
	SortBy         string // best_price, newest, completion_rate, trade_count
	Limit          int
	Offset         int
}
