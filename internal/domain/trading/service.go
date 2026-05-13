// Package trading provides the core P2P trading domain service.
//
// Implementations are split across focused sub-files:
//   - ad.go    – ad management (ListActiveAds, GetAd, CreateAd, UpdateAd, DeleteAd, ListUserAds, ToggleAdStatus)
//   - trade.go – trade lifecycle (InitiateTrade, ListTrades, ListAllTrades, CountTrades, GetTrade,
//     MarkAsPaid, ReleaseEscrow, CancelTrade, DisputeTrade,
//     ReconcileExpiredTrades, FlagAutoDisputesForOverduePaidTrades)
//   - chat.go  – chat & feedback (SendMessage, SendFileMessage, GetChatHistory, LeaveFeedback)
package trading

import (
	"context"
	"cryplio/internal/domain/dispute"
	"cryplio/internal/domain/identity"
	"cryplio/internal/domain/notification"
	"cryplio/internal/domain/platform"
	"cryplio/pkg/config"
	"time"

	"github.com/google/uuid"
)

// TradeService is the primary service interface for the trading domain.
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

	// Ad Management
	ListUserAds(ctx context.Context, userID uuid.UUID) ([]TradeAd, int, error)
	ToggleAdStatus(ctx context.Context, adID, userID uuid.UUID) error
}

type tradeService struct {
	tradeRepo           TradeRepository
	identityRepo        identity.UserRepository
	disputeRepo         dispute.Repository
	escrowClient        EscrowContractClient
	notificationService notification.Service
	platformRepo        platform.PlatformRepository
	cfg                 *config.Config
}

// NewTradeService constructs a TradeService with all required dependencies.
func NewTradeService(
	tradeRepo TradeRepository,
	identityRepo identity.UserRepository,
	disputeRepo dispute.Repository,
	escrowClient EscrowContractClient,
	notificationService notification.Service,
	platformRepo platform.PlatformRepository,
	cfg *config.Config,
) TradeService {
	return &tradeService{
		tradeRepo:           tradeRepo,
		identityRepo:        identityRepo,
		disputeRepo:         disputeRepo,
		escrowClient:        escrowClient,
		notificationService: notificationService,
		platformRepo:        platformRepo,
		cfg:                 cfg,
	}
}
