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

	// Blockchain Clients (Mandatory for production)
	if cfg.EthRPCURL == "" || cfg.EthPrivateKey == "" {
		return nil, fmt.Errorf("blockchain configuration missing: ETH_RPC_URL and ETH_PRIVATE_KEY are required")
	}

	walletClient, err := blockchain.NewEvmWalletClient(cfg.EthRPCURL, cfg.EthPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("initialize EVM Wallet Client: %w", err)
	}

	var escrowClient domaintrading.EscrowContractClient
	if cfg.EscrowContractAddress != "" && cfg.EscrowABIPath != "" {
		escrowClient, err = blockchain.NewEvmEscrowClient(cfg.EthRPCURL, cfg.EthPrivateKey, cfg.EscrowContractAddress, cfg.EscrowABIPath)
		if err != nil {
			return nil, fmt.Errorf("initialize EVM Escrow Client: %w", err)
		}
	} else {
		return nil, fmt.Errorf("escrow configuration missing: ESCROW_CONTRACT_ADDRESS and ESCROW_ABI_PATH are required")
	}

	// Create wallet service first so it can be injected into auth service
	walletRepo := walletpostgres.NewWalletRepository(db)
	platformRepo := platformpostgres.NewPlatformRepository(db)
	platformService := platform.NewPlatformService(platformRepo)
	walletService := domainwallet.NewService(walletRepo, walletClient, platformService, cfg)

	// Create auth service with wallet service for auto-creating wallets on registration
	authService := domainidentity.NewAuthService(domainidentity.AuthServiceConfig{
		UserRepo:           userRepo,
		JWTSecret:          cfg.JWTSecret,
		JWTExpiry:          cfg.JWTExpiry,
		RefreshTokenExpiry: cfg.RefreshTokenExpiry,
		CookieName:         cfg.CookieName,
		CookieSecure:       cfg.CookieSecure,
		CookieSameSite:     cfg.CookieSameSite,
		IssuerName:         cfg.IssuerName,
	}).WithGoogleOAuth(
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
	tradeService := domaintrading.NewTradeService(domaintrading.TradeServiceConfig{
		TradeRepo:           tradeRepo,
		IdentityRepo:        userRepo,
		DisputeRepo:         disputeRepo,
		EscrowClient:        escrowClient,
		NotificationService: notificationService,
		PlatformRepo:        platformRepo,
		WalletRepo:          walletRepo,
		Cfg:                 cfg,
	})

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
