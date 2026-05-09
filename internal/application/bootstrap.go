package application

import (
	"context"
	"database/sql"
	"fmt"

	domaindispute "cryplio/internal/domain/dispute"
	domainidentity "cryplio/internal/domain/identity"
	"cryplio/internal/domain/market"
	domainnotification "cryplio/internal/domain/notification"
	"cryplio/internal/domain/platform"
	domaintrading "cryplio/internal/domain/trading"
	domainwallet "cryplio/internal/domain/wallet"
	"cryplio/internal/infrastructure/blockchain"
	"cryplio/internal/infrastructure/notification"
	disputepostgres "cryplio/internal/infrastructure/persistence/postgres/dispute"
	identitypostgres "cryplio/internal/infrastructure/persistence/postgres/identity"
	notificationpostgres "cryplio/internal/infrastructure/persistence/postgres/notification"
	platformpostgres "cryplio/internal/infrastructure/persistence/postgres/platform"
	tradingpostgres "cryplio/internal/infrastructure/persistence/postgres/trading"
	walletpostgres "cryplio/internal/infrastructure/persistence/postgres/wallet"
	"cryplio/internal/infrastructure/persistence/redis"
	"cryplio/internal/infrastructure/storage"
	"cryplio/internal/infrastructure/worker"
	httpapi "cryplio/internal/interfaces/http"
	"cryplio/internal/interfaces/websocket"
	"cryplio/pkg/config"
	"cryplio/pkg/database"
	"cryplio/pkg/logger"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config    *config.Config
	DB        *sql.DB
	Router    *gin.Engine
	Worker    *worker.Worker
	Scheduler *worker.Scheduler
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	db, err := database.Open(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	dbPingFn := func(ctx context.Context) error { return db.PingContext(ctx) }

	userRepo := identitypostgres.NewUserRepository(db)
	emailClient := notification.NewSMTPClient(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.EmailFrom,
		cfg.FrontendURL,
	)

	// Blockchain Clients (create first so wallet service can be injected into auth service)
	var escrowClient domaintrading.EscrowContractClient
	var walletClient domainwallet.WalletClient

	if cfg.EscrowContractAddress != "" && cfg.EthRPCURL != "" {
		evmEscrow, err := blockchain.NewEvmEscrowClient(cfg.EthRPCURL, cfg.EthPrivateKey, cfg.EscrowContractAddress)
		if err == nil {
			escrowClient = evmEscrow
		} else {
			logger.Warn("failed to init EVM Escrow Client", logger.Fields{"error": err.Error()})
			escrowClient = blockchain.NewMockEscrowContractClient()
		}
	} else {
		// Use mock escrow service for MVP without smart contracts
		escrowClient = blockchain.NewMockEscrowContractClient()
	}

	if cfg.EthPrivateKey != "" {
		evmWallet, err := blockchain.NewEvmWalletClient(cfg.EthRPCURL, cfg.EthPrivateKey)
		if err == nil {
			walletClient = evmWallet
		} else {
			logger.Warn("failed to init EVM Wallet Client", logger.Fields{"error": err.Error()})
			walletClient = blockchain.NewNoopWalletClient()
		}
	} else {
		escrowClient = blockchain.NewNoopEscrowContractClient()
		walletClient = blockchain.NewNoopWalletClient()
	}

	// Create wallet service first so it can be injected into auth service
	walletRepo := walletpostgres.NewWalletRepository(db)
	walletService := domainwallet.NewService(walletRepo, walletClient)

	// Create auth service with wallet service for auto-creating wallets on registration
	authService := domainidentity.NewAuthService(
		userRepo,
		cfg.JWTSecret,
		cfg.JWTExpiry,
		cfg.RefreshTokenExpiry,
		cfg.CookieName,
		cfg.CookieSecure,
		cfg.CookieSameSite,
		cfg.IssuerName,
	).WithGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.OAuthRedirectURL,
	).WithPasswordResetMailer(emailClient).
		WithWalletService(walletService)

	notificationRepo := notificationpostgres.NewNotificationRepository(db)
	notificationService := domainnotification.NewService(notificationRepo, emailClient)

	disputeRepo := disputepostgres.NewDisputeRepository(db)
	disputeService := domaindispute.NewService(disputeRepo)

	tradeRepo := tradingpostgres.NewTradeRepository(db)
	platformRepo := platformpostgres.NewPlatformRepository(db)
	tradeService := domaintrading.NewTradeService(tradeRepo, userRepo, disputeRepo, escrowClient, notificationService, platformRepo, cfg)
	platformService := platform.NewPlatformService(platformRepo)

	rateService := market.NewRateService()

	storage, err := storage.NewS3Storage(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize storage: %w", err)
	}

	// Create WebSocket service
	wsService := websocket.NewService()

	// Create WebSocket notifier for real-time notifications
	emailService := notification.NewMockEmailService("noreply@cryplio.com", "Cryplio")
	wsNotifier := notification.NewWebSocketNotifier(wsService, emailService)

	// Connect WebSocket notifier to notification service for real-time delivery
	notificationService.SetWebSocketNotifier(wsNotifier)

	// Redis rate limiter — falls back to in-memory if Redis is unavailable.
	var rateLimiter redis.RateLimiter
	redisClient, redisErr := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if redisErr != nil {
		logger.Warn("Redis unavailable; using in-memory rate limiter (not suitable for multi-instance deployments)", logger.Fields{"error": redisErr.Error()})
	} else {
		rateLimiter = redis.NewRedisRateLimiter(redisClient)
		logger.Info("Redis rate limiter active", logger.Fields{"addr": cfg.RedisAddr})
	}

	router := httpapi.SetupRouter(cfg, authService, tradeService, platformService, walletService, disputeService, notificationService, storage, rateService, wsService, rateLimiter, dbPingFn)

	asynqWorker := worker.NewWorker(cfg, tradeService)
	asynqScheduler := worker.NewScheduler(cfg)

	return &App{
		Config:    cfg,
		DB:        db,
		Router:    router,
		Worker:    asynqWorker,
		Scheduler: asynqScheduler,
	}, nil
}
